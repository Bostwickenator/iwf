package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type ContinueAsNewer struct {
	pendingStateExecutionsRequestCommands   map[string]service.PendingStateExecutionRequestCommands
	pendingStateExecutionsCompletedCommands map[string]service.PendingStateExecutionCompletedCommands
	interStateChannel                       *InterStateChannel
	stateExecutionCounter                   *StateExecutionCounter
	persistenceManager                      *PersistenceManager
}

func NewContinueAsNewer(
	interStateChannel *InterStateChannel, stateExecutionCounter *StateExecutionCounter, persistenceManager *PersistenceManager,
) *ContinueAsNewer {
	return &ContinueAsNewer{
		interStateChannel:                       interStateChannel,
		stateExecutionCounter:                   stateExecutionCounter,
		persistenceManager:                      persistenceManager,
		pendingStateExecutionsCompletedCommands: map[string]service.PendingStateExecutionCompletedCommands{},
		pendingStateExecutionsRequestCommands:   map[string]service.PendingStateExecutionRequestCommands{},
	}
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext, provider WorkflowProvider) error {
	err := provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func() (*service.DumpAllInternalResponse, error) {
		return &service.DumpAllInternalResponse{
			InterStateChannelReceived:               c.interStateChannel.ReadReceived(nil),
			StateExecutionCounterInfo:               c.stateExecutionCounter.Dump(),
			PendingStateExecutionsCompletedCommands: c.pendingStateExecutionsCompletedCommands,
			PendingStateExecutionsRequestCommands:   c.pendingStateExecutionsRequestCommands,
			DataObjects:                             c.persistenceManager.GetAllDataObjects(),
			SearchAttributes:                        c.persistenceManager.GetAllSearchAttributes(),
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ContinueAsNewer) AddPendingStateExecution(
	stateExecutionId string,
	completedTimerCommands map[int]bool, completedSignalCommands, completedInterStateChannelCommands map[int]*iwfidl.EncodedObject,
	timerCommands []iwfidl.TimerCommand, signalCommands []iwfidl.SignalCommand, interStateChannelCommands []iwfidl.InterStateChannelCommand,
) {
	c.pendingStateExecutionsCompletedCommands[stateExecutionId] = service.PendingStateExecutionCompletedCommands{
		CompletedTimerCommands:             completedTimerCommands,
		CompletedSignalCommands:            completedSignalCommands,
		CompletedInterStateChannelCommands: completedInterStateChannelCommands,
	}
	c.pendingStateExecutionsRequestCommands[stateExecutionId] = service.PendingStateExecutionRequestCommands{
		TimerCommands:             timerCommands,
		SignalCommands:            signalCommands,
		InterStateChannelCommands: interStateChannelCommands,
	}
}

func (c *ContinueAsNewer) DeletePendingStateExecution(stateExecutionId string) {
	delete(c.pendingStateExecutionsCompletedCommands, stateExecutionId)
	delete(c.pendingStateExecutionsRequestCommands, stateExecutionId)
}
