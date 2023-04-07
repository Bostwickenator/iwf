package interpreter

import (
	"fmt"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

func InterpreterImpl(ctx UnifiedContext, provider WorkflowProvider, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	var err error
	globalVersioner := NewGlobalVersioner(provider, ctx)
	if globalVersioner.IsAfterVersionOfUsingGlobalVersioning() {
		err = globalVersioner.UpsertGlobalVersionSearchAttribute()
		if err != nil {
			return nil, err
		}
	}

	if !input.Config.GetDisableSystemSearchAttribute() {
		if !globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() {
			// we have stopped upsert here in new versions, because it's done in start workflow request
			err = provider.UpsertSearchAttributes(ctx, map[string]interface{}{
				service.SearchAttributeIwfWorkflowType: input.IwfWorkflowType,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	workflowConfiger := NewWorkflowConfiger(input.Config)
	basicInfo := service.BasicInfo{
		IwfWorkflowType: input.IwfWorkflowType,
		IwfWorkerUrl:    input.IwfWorkerUrl,
	}

	var interStateChannel *InterStateChannel
	var stateRequestQueue *StateRequestQueue
	var persistenceManager *PersistenceManager
	var timerProcessor *TimerProcessor
	var continueAsNewCounter *ContinueAsNewCounter
	var signalReceiver *SignalReceiver
	var stateExecutionCounter *StateExecutionCounter
	var outputCollector *OutputCollector
	var continueAsNewer *ContinueAsNewer
	if input.IsResumeFromContinueAsNew {
		canInput := input.ContinueAsNewInput
		config := workflowConfiger.Get()
		previous, err := LoadInternalsFromPreviousRun(ctx, provider, canInput.PreviousInternalRunId, config.GetContinueAsNewPageSizeInBytes())
		if err != nil {
			return nil, err
		}

		// The below initialization order should be the same as for non-continueAsNew

		interStateChannel = RebuildInterStateChannel(previous.InterStateChannelReceived)
		stateRequestQueue = NewStateRequestQueueWithResumeRequests(previous.StatesToStartFromBeginning, previous.StateExecutionsToResume)
		persistenceManager = RebuildPersistenceManager(provider, previous.DataObjects, previous.SearchAttributes)
		timerProcessor = NewTimerProcessor(ctx, provider, previous.StaleSkipTimerSignals)
		continueAsNewCounter = NewContinueAsCounter(workflowConfiger, ctx, provider)
		signalReceiver = NewSignalReceiver(ctx, provider, timerProcessor, continueAsNewCounter, workflowConfiger, previous.SignalsReceived)
		counterInfo := previous.StateExecutionCounterInfo
		stateExecutionCounter = RebuildStateExecutionCounter(ctx, provider,
			counterInfo.StateIdStartedCount, counterInfo.StateIdCurrentlyExecutingCount, counterInfo.TotalCurrentlyExecutingCount,
			workflowConfiger, continueAsNewCounter)
		outputCollector = NewOutputCollector(previous.StateOutputs)
		continueAsNewer = NewContinueAsNewer(provider, interStateChannel, signalReceiver, stateExecutionCounter, persistenceManager, stateRequestQueue, outputCollector, timerProcessor)
	} else {
		interStateChannel = NewInterStateChannel()

		stateRequestQueue = NewStateRequestQueue(iwfidl.StateMovement{
			StateId:      input.StartStateId,
			StateOptions: &input.StateOptions,
			StateInput:   &input.StateInput,
		})
		persistenceManager = NewPersistenceManager(provider, input.InitSearchAttributes)
		timerProcessor = NewTimerProcessor(ctx, provider, nil)
		continueAsNewCounter = NewContinueAsCounter(workflowConfiger, ctx, provider)
		signalReceiver = NewSignalReceiver(ctx, provider, timerProcessor, continueAsNewCounter, workflowConfiger, nil)
		stateExecutionCounter = NewStateExecutionCounter(ctx, provider, workflowConfiger, continueAsNewCounter)
		outputCollector = NewOutputCollector(nil)
		continueAsNewer = NewContinueAsNewer(provider, interStateChannel, signalReceiver, stateExecutionCounter, persistenceManager, stateRequestQueue, outputCollector, timerProcessor)
	}

	err = provider.SetQueryHandler(ctx, service.GetDataObjectsWorkflowQueryType, func(req service.GetDataObjectsQueryRequest) (service.GetDataObjectsQueryResponse, error) {
		return persistenceManager.GetDataObjectsByKey(req), nil
	})
	if err != nil {
		return nil, err
	}
	err = provider.SetQueryHandler(ctx, service.GetSearchAttributesWorkflowQueryType, func() ([]iwfidl.SearchAttribute, error) {
		return persistenceManager.GetAllSearchAttributes(), nil
	})
	if err != nil {
		return nil, err
	}
	err = continueAsNewer.SetQueryHandlersForContinueAsNew(ctx)
	if err != nil {
		return nil, err
	}
	err = provider.SetQueryHandler(ctx, service.DebugDumpQueryType, func() (*service.DebugDumpResponse, error) {
		return &service.DebugDumpResponse{
			Config: workflowConfiger.Get(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	var errToFailWf error // Note that today different errors could overwrite each other, we only support last one wins. we may use multiError to improve.
	var forceCompleteWf bool

	// this is for an optimization for StateId Search attribute, see updateStateIdSearchAttribute in stateExecutionCounter
	defer stateExecutionCounter.ClearExecutingStateIdsSearchAttributeFinally()

	for !stateRequestQueue.IsEmpty() {

		statesToExecute := stateRequestQueue.TakeAll()
		err = stateExecutionCounter.MarkStateIdExecutingIfNotYet(statesToExecute)
		if err != nil {
			return nil, err
		}

		for _, stateReqForLoopingOnly := range statesToExecute {
			// execute in another thread for parallelism
			// state must be passed via parameter https://stackoverflow.com/questions/67263092
			stateCtx := provider.ExtendContextWithValue(ctx, "stateReq", stateReqForLoopingOnly)
			provider.GoNamed(stateCtx, "state-execution-thread:"+stateReqForLoopingOnly.GetStateId(), func(ctx UnifiedContext) {
				stateReq, ok := provider.GetContextValue(ctx, "stateReq").(StateRequest)
				if !ok {
					errToFailWf = provider.NewApplicationError(
						string(iwfidl.SERVER_INTERNAL_ERROR_TYPE),
						"critical code bug when passing state request via context",
					)
					return
				}

				var state iwfidl.StateMovement
				var stateExeId string
				if stateReq.IsResumeRequest() {
					resumeReq := stateReq.GetStateResumeRequest()
					state = resumeReq.State
					stateExeId = resumeReq.StateExecutionId
				} else {
					state = stateReq.GetStateStartRequest()
					stateExeId = stateExecutionCounter.CreateNextExecutionId(state.GetStateId())
				}

				decision, stateExecStatus, err := executeState(
					ctx, provider, basicInfo, stateReq, stateExeId, persistenceManager, interStateChannel,
					signalReceiver, timerProcessor, continueAsNewer, continueAsNewCounter)
				if err != nil {
					errToFailWf = err
					// state execution fail should fail the workflow, no more processing
					return
				}

				if stateExecStatus == service.CompletedStateExecutionStatus {
					// NOTE: decision is only available on this CompletedStateExecutionStatus

					shouldClose, gracefulComplete, forceComplete, forceFail, output, err := checkClosingWorkflow(provider, decision, state.GetStateId(), stateExeId)
					if err != nil {
						errToFailWf = err
						// no return so that it can fall through to call MarkStateExecutionCompleted
					}
					if gracefulComplete || forceComplete || forceFail {
						outputCollector.Add(*output)
					}
					if forceComplete {
						forceCompleteWf = true
					}
					if forceFail {
						errToFailWf = provider.NewApplicationError(
							string(iwfidl.STATE_DECISION_FAILING_WORKFLOW_ERROR_TYPE),
							outputCollector.GetAll(),
						)
						// no return so that it can fall through to call MarkStateExecutionCompleted
					}
					if !shouldClose && decision.HasNextStates() {
						stateRequestQueue.AddStateStartRequests(decision.GetNextStates())
					}

					// finally, mark state completed and may also update system search attribute(IwfExecutingStateIds)
					err = stateExecutionCounter.MarkStateExecutionCompleted(state)
					if err != nil {
						errToFailWf = err
					}
				}
			}) // end of executing one state
		} // end loop of executing all states from the queue for one iteration

		// The conditions here are quite tricky:
		// For !stateRequestQueue.IsEmpty(): We need some condition to wait here because all the state execution are running in different thread.
		//    Right after the queue are popped it becomes empty. When it's not empty, it means there are new states to execute pushed into the queue,
		//    and it's time to wake up the outer loop to go to next iteration. Alternatively, waiting for all current started in this iteration to complete will also work,
		//    but not as efficient as this one because it will take much longer time.
		// For errToFailWf != nil || forceCompleteWf: this means we need to close workflow immediately
		// For stateExecutionCounter.GetTotalCurrentlyExecutingCount() == 0: this means all the state executions have reach "Dead Ends" so the workflow can complete gracefully without output
		// For continueAsNewCounter.IsThresholdMet(): this means workflow need to continueAsNew
		awaitError := provider.Await(ctx, func() bool {
			failByApi, errStr := signalReceiver.IsFailWorkflowRequested()
			if failByApi {
				errToFailWf = provider.NewApplicationError(
					string(iwfidl.CLIENT_API_FAILING_WORKFLOW_ERROR_TYPE),
					errStr,
				)

				return true
			}
			return !stateRequestQueue.IsEmpty() || errToFailWf != nil || forceCompleteWf || stateExecutionCounter.GetTotalCurrentlyExecutingCount() == 0 || continueAsNewCounter.IsThresholdMet()
		})
		if continueAsNewCounter.IsThresholdMet() {
			// NOTE: drain thread before checking errToFailWf/forceCompleteWf so that we can close the workflow if possible
			err := continueAsNewer.DrainThreads(ctx)
			if err != nil {
				awaitError = err
			}
		}

		if errToFailWf != nil || forceCompleteWf {
			return &service.InterpreterWorkflowOutput{
				StateCompletionOutputs: outputCollector.GetAll(),
			}, errToFailWf
		}

		if awaitError != nil {
			// this could happen for cancellation
			errToFailWf = awaitError
			break
		}
		if stateRequestQueue.IsEmpty() && !continueAsNewer.HasAnyStateExecutionToResume() {
			// if it is empty and no stateExecutionsToResume, then we don't need to do continue as new -- the workflow has reached "dead ends" and will just complete after the loop
			break
		}
		if continueAsNewCounter.IsThresholdMet() {
			// NOTE: This must be the last thing before continueAsNew!!!
			// Otherwise, there could be signals unhandled
			signalReceiver.DrainAllUnreceivedSignals(ctx)

			// at here, all signals + threads are drained, so it's safe to continueAsNew
			// after draining signals, there could be some changes
			// 1. last fail workflow signal
			failByApi, errStr := signalReceiver.IsFailWorkflowRequested()
			if failByApi {
				return &service.InterpreterWorkflowOutput{
						StateCompletionOutputs: outputCollector.GetAll(),
					}, provider.NewApplicationError(
						string(iwfidl.CLIENT_API_FAILING_WORKFLOW_ERROR_TYPE),
						errStr,
					)
			}
			// 2. last update config, do it here because we use input to carry over config, not continueAsNewer query
			input.Config = workflowConfiger.Get() // update config to the lastest before continueAsNew to carry over
			input.IsResumeFromContinueAsNew = true
			input.ContinueAsNewInput = service.ContinueAsNewInput{
				PreviousInternalRunId: provider.GetWorkflowInfo(ctx).WorkflowExecution.RunID,
			}
			return nil, provider.NewInterpreterContinueAsNewError(ctx, input)
		}
	} // end main loop -- loop until no more state can be executed (dead end)

	// gracefully complete workflow when all states are executed to dead ends
	return &service.InterpreterWorkflowOutput{
		StateCompletionOutputs: outputCollector.GetAll(),
	}, errToFailWf
}

func checkClosingWorkflow(
	provider WorkflowProvider, decision *iwfidl.StateDecision, currentStateId, currentStateExeId string,
) (shouldClose, gracefulComplete, forceComplete, forceFail bool, completeOutput *iwfidl.StateCompletionOutput, err error) {
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.GracefulCompletingWorkflowStateId {
			shouldClose = true
			gracefulComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceCompletingWorkflowStateId {
			shouldClose = true
			forceComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceFailingWorkflowStateId {
			shouldClose = true
			forceFail = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
	}
	if shouldClose && len(decision.NextStates) > 1 {
		// Illegal decision
		err = provider.NewApplicationError(
			string(iwfidl.INVALID_USER_WORKFLOW_CODE_ERROR_TYPE),
			"invalid state decisions. Closing workflow decision cannot be combined with other state decisions",
		)
		return
	}
	return
}

func executeState(
	ctx UnifiedContext,
	provider WorkflowProvider,
	basicInfo service.BasicInfo,
	stateReq StateRequest,
	stateExeId string,
	persistenceManager *PersistenceManager,
	interStateChannel *InterStateChannel,
	signalReceiver *SignalReceiver,
	timerProcessor *TimerProcessor,
	continueAsNewer *ContinueAsNewer,
	continueAsNewCounter *ContinueAsNewCounter,
) (*iwfidl.StateDecision, service.StateExecutionStatus, error) {
	info := provider.GetWorkflowInfo(ctx)
	executionContext := iwfidl.Context{
		WorkflowId:               info.WorkflowExecution.ID,
		WorkflowRunId:            info.WorkflowExecution.RunID,
		WorkflowStartedTimestamp: info.WorkflowStartTime.Unix(),
		StateExecutionId:         stateExeId,
	}
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}

	var err error
	var errStartApi error
	var startResponse *iwfidl.WorkflowStateStartResponse
	var stateExecutionLocal []iwfidl.KeyValue
	var commandReq iwfidl.CommandRequest
	commandReqDoneOrCanceled := false
	completedTimerCmds := map[int]service.InternalTimerStatus{}
	completedSignalCmds := map[int]*iwfidl.EncodedObject{}
	completedInterStateChannelCmds := map[int]*iwfidl.EncodedObject{}

	var state iwfidl.StateMovement
	isResumeFromContinueAsNew := stateReq.IsResumeRequest()
	if isResumeFromContinueAsNew {
		state = stateReq.GetStateResumeRequest().State
	} else {
		state = stateReq.GetStateStartRequest()
	}

	if isResumeFromContinueAsNew {
		resumeStateRequest := stateReq.GetStateResumeRequest()
		stateExecutionLocal = resumeStateRequest.StateExecutionLocals
		commandReq = resumeStateRequest.CommandRequest
		completedCmds := resumeStateRequest.StateExecutionCompletedCommands
		completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds = completedCmds.CompletedTimerCommands, completedCmds.CompletedSignalCommands, completedCmds.CompletedInterStateChannelCommands
	} else {
		if state.StateOptions != nil {
			if state.StateOptions.GetStartApiTimeoutSeconds() > 0 {
				activityOptions.StartToCloseTimeout = time.Duration(state.StateOptions.GetStartApiTimeoutSeconds()) * time.Second
			}
			activityOptions.RetryPolicy = state.StateOptions.StartApiRetryPolicy
		}

		ctx = provider.WithActivityOptions(ctx, activityOptions)

		errStartApi = provider.ExecuteActivity(ctx, StateStart, provider.GetBackendType(), service.StateStartActivityInput{
			IwfWorkerUrl: basicInfo.IwfWorkerUrl,
			Request: iwfidl.WorkflowStateStartRequest{
				Context:          executionContext,
				WorkflowType:     basicInfo.IwfWorkflowType,
				WorkflowStateId:  state.StateId,
				StateInput:       state.StateInput,
				SearchAttributes: persistenceManager.LoadSearchAttributes(state.StateOptions),
				DataObjects:      persistenceManager.LoadDataObjects(state.StateOptions),
			},
		}).Get(ctx, &startResponse)

		if errStartApi != nil && !shouldProceedOnStartApiError(state) {
			return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, errStartApi)
		}

		err := persistenceManager.ProcessUpsertSearchAttribute(ctx, startResponse.GetUpsertSearchAttributes())
		if err != nil {
			return nil, service.FailureStateExecutionStatus, err
		}
		err = persistenceManager.ProcessUpsertDataObject(startResponse.GetUpsertDataObjects())
		if err != nil {
			return nil, service.FailureStateExecutionStatus, err
		}
		interStateChannel.ProcessPublishing(startResponse.GetPublishToInterStateChannel())

		commandReq = startResponse.GetCommandRequest()
		stateExecutionLocal = startResponse.GetUpsertStateLocals()
	}

	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.AddTimers(stateExeId, commandReq.GetTimerCommands(), completedTimerCmds)
		for idx, cmd := range commandReq.GetTimerCommands() {
			if _, ok := completedTimerCmds[idx]; ok {
				// skip the completed timers(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("timer", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				// Note that commandReqDoneOrCanceled is needed for two cases:
				// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
				// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
				status := timerProcessor.WaitForTimerFiredOrSkipped(ctx, stateExeId, idx, &commandReqDoneOrCanceled)
				if status == service.TimerSkipped || status == service.TimerFired {
					completedTimerCmds[idx] = status
				}
			})
		}
	}

	if len(commandReq.GetSignalCommands()) > 0 {
		for idx, cmd := range commandReq.GetSignalCommands() {
			if _, ok := completedSignalCmds[idx]; ok {
				// skip completed signal(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("signal", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				cmd, ok := provider.GetContextValue(ctx, "cmd").(iwfidl.SignalCommand)
				if !ok {
					panic("critical code bug")
				}
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}
				received := false
				_ = provider.Await(ctx, func() bool {
					received = signalReceiver.HasSignal(cmd.SignalChannelName)
					// Note that commandReqDoneOrCanceled is needed for two cases:
					// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
					// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
					return received || commandReqDoneOrCanceled
				})
				if received {
					completedSignalCmds[idx] = signalReceiver.Retrieve(cmd.SignalChannelName)
				}
			})
		}
	}

	if len(commandReq.GetInterStateChannelCommands()) > 0 {
		for idx, cmd := range commandReq.GetInterStateChannelCommands() {
			if _, ok := completedInterStateChannelCmds[idx]; ok {
				// skip completed interStateChannelCommand(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("interstate", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				cmd, ok := provider.GetContextValue(ctx, "cmd").(iwfidl.InterStateChannelCommand)
				if !ok {
					panic("critical code bug")
				}
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				received := false
				_ = provider.Await(ctx, func() bool {
					received = interStateChannel.HasData(cmd.ChannelName)
					// Note that commandReqDoneOrCanceled is needed for two cases:
					// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
					// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
					return received || commandReqDoneOrCanceled
				})

				if received {
					completedInterStateChannelCmds[idx] = interStateChannel.Retrieve(cmd.ChannelName)
				}
			})
		}
	}

	continueAsNewer.AddPotentialStateExecutionToResume(
		stateExeId, state, stateExecutionLocal, commandReq,
		completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds,
	)
	_ = provider.Await(ctx, func() bool {
		return IsDeciderTriggerConditionMet(commandReq, completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds) || continueAsNewCounter.IsThresholdMet()
	})
	commandReqDoneOrCanceled = true
	if !IsDeciderTriggerConditionMet(commandReq, completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds) {
		// this means continueAsNewCounter.IsThresholdMet == true
		// not using continueAsNewCounter.IsThresholdMet because deciderTrigger is higher prioritized
		// it won't continueAsNew in those cases 1. start Api fail with proceed policy, 2. empty commands, 3. both commands and continueAsNew are met
		return nil, service.WaitingCommandsStateExecutionStatus, nil
	}

	commandRes := &iwfidl.CommandResults{}
	commandRes.StateStartApiSucceeded = iwfidl.PtrBool(errStartApi == nil)

	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.RemovePendingTimersOfState(stateExeId)

		var timerResults []iwfidl.TimerResult
		for idx, cmd := range commandReq.GetTimerCommands() {
			status := iwfidl.SCHEDULED
			if _, ok := completedTimerCmds[idx]; ok {
				// TODO expose skipped status to external
				status = iwfidl.FIRED
			}
			timerResults = append(timerResults, iwfidl.TimerResult{
				CommandId:   cmd.GetCommandId(),
				TimerStatus: status,
			})
		}
		commandRes.SetTimerResults(timerResults)
	}

	if len(commandReq.GetSignalCommands()) > 0 {
		var signalResults []iwfidl.SignalResult
		for idx, cmd := range commandReq.GetSignalCommands() {
			status := iwfidl.RECEIVED
			result, completed := completedSignalCmds[idx]
			if !completed {
				status = iwfidl.WAITING
			}

			signalResults = append(signalResults, iwfidl.SignalResult{
				CommandId:           cmd.GetCommandId(),
				SignalChannelName:   cmd.GetSignalChannelName(),
				SignalValue:         result,
				SignalRequestStatus: status,
			})
		}
		commandRes.SetSignalResults(signalResults)
	}

	if len(commandReq.GetInterStateChannelCommands()) > 0 {
		var interStateChannelResults []iwfidl.InterStateChannelResult
		for idx, cmd := range commandReq.GetInterStateChannelCommands() {
			status := iwfidl.RECEIVED
			result, completed := completedInterStateChannelCmds[idx]
			if !completed {
				status = iwfidl.WAITING
			}

			interStateChannelResults = append(interStateChannelResults, iwfidl.InterStateChannelResult{
				CommandId:     cmd.CommandId,
				ChannelName:   cmd.ChannelName,
				RequestStatus: status,
				Value:         result,
			})
		}
		commandRes.SetInterStateChannelResults(interStateChannelResults)
	}

	activityOptions = ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	if state.StateOptions != nil {
		if state.StateOptions.GetDecideApiTimeoutSeconds() > 0 {
			activityOptions.StartToCloseTimeout = time.Duration(state.StateOptions.GetDecideApiTimeoutSeconds()) * time.Second
		}
		activityOptions.RetryPolicy = state.StateOptions.DecideApiRetryPolicy
	}

	ctx = provider.WithActivityOptions(ctx, activityOptions)
	var decideResponse *iwfidl.WorkflowStateDecideResponse
	err = provider.ExecuteActivity(ctx, StateDecide, provider.GetBackendType(), service.StateDecideActivityInput{
		IwfWorkerUrl: basicInfo.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateDecideRequest{
			Context:          executionContext,
			WorkflowType:     basicInfo.IwfWorkflowType,
			WorkflowStateId:  state.StateId,
			CommandResults:   commandRes,
			StateLocals:      stateExecutionLocal,
			SearchAttributes: persistenceManager.LoadSearchAttributes(state.StateOptions),
			DataObjects:      persistenceManager.LoadDataObjects(state.StateOptions),
			StateInput:       state.StateInput,
		},
	}).Get(ctx, &decideResponse)
	if err != nil {
		return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, err)
	}

	err = persistenceManager.ProcessUpsertSearchAttribute(ctx, decideResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	err = persistenceManager.ProcessUpsertDataObject(decideResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	interStateChannel.ProcessPublishing(decideResponse.GetPublishToInterStateChannel())

	continueAsNewer.RemoveStateExecutionToResume(stateExeId)

	decision := decideResponse.GetStateDecision()
	return &decision, service.CompletedStateExecutionStatus, nil
}

func shouldProceedOnStartApiError(state iwfidl.StateMovement) bool {
	if state.StateOptions == nil {
		return false
	}

	if state.StateOptions.StartApiFailurePolicy == nil {
		return false
	}

	return state.StateOptions.GetStartApiFailurePolicy() == iwfidl.PROCEED_TO_DECIDE_ON_START_API_FAILURE
}

func convertStateApiActivityError(provider WorkflowProvider, err error) error {
	if provider.IsApplicationError(err) {
		return err
	}
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE), err.Error())
}

func getCommandThreadName(prefix string, stateExecId, cmdId string, idx int) string {
	return fmt.Sprintf("%v-%v-%v-%v", prefix, stateExecId, cmdId, idx)
}
