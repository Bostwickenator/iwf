package api

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"log"
	"net/http"
	"time"
)

type serviceImpl struct {
	client UnifiedClient
}

func (s *serviceImpl) Close() {
	s.client.Close()
}

func NewApiService(client UnifiedClient) (ApiService, error) {
	return &serviceImpl{
		client: client,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowStartPost(req iwfidl.WorkflowStartRequest) (*iwfidl.WorkflowStartResponse, *ErrorAndStatus) {
	workflowOptions := StartWorkflowOptions{
		ID:                 req.GetWorkflowId(),
		TaskQueue:          service.TaskQueue,
		WorkflowRunTimeout: time.Duration(req.WorkflowTimeoutSeconds) * time.Second,
	}

	input := service.InterpreterWorkflowInput{
		IwfWorkflowType: req.GetIwfWorkflowType(),
		IwfWorkerUrl:    req.GetIwfWorkerUrl(),
		StartStateId:    req.GetStartStateId(),
		StateInput:      req.GetStateInput(),
		StateOptions:    req.GetStateOptions(),
	}
	runId, err := s.client.StartInterpreterWorkflow(context.Background(), workflowOptions, input)
	if err != nil {
		return nil, s.handleError(err)

	}

	log.Println("Started workflow", "WorkflowID", req.WorkflowId, "RunID", runId)

	return &iwfidl.WorkflowStartResponse{
		WorkflowRunId: iwfidl.PtrString(runId),
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSignalPost(req iwfidl.WorkflowSignalRequest) *ErrorAndStatus {
	err := s.client.SignalWorkflow(context.Background(),
		req.GetWorkflowId(), req.GetWorkflowRunId(), req.GetSignalChannelName(), req.GetSignalValue())
	if err != nil {
		return s.handleError(err)
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetQueryAttributesPost(req iwfidl.WorkflowGetQueryAttributesRequest) (*iwfidl.WorkflowGetQueryAttributesResponse, *ErrorAndStatus) {
	var queryResult1 service.QueryAttributeResponse
	err := s.client.QueryWorkflow(context.Background(), &queryResult1,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.AttributeQueryType,
		service.QueryAttributeRequest{
			Keys: req.AttributeKeys,
		})

	if err != nil {
		return nil, s.handleError(err)
	}

	return &iwfidl.WorkflowGetQueryAttributesResponse{
		QueryAttributes: queryResult1.AttributeValues,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowGetPost(req iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus) {
	return s.doApiV1WorkflowGetPost(req, false)
}

func (s *serviceImpl) ApiV1WorkflowGetWithWaitPost(req iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus) {
	return s.doApiV1WorkflowGetPost(req, true)
}

func (s *serviceImpl) doApiV1WorkflowGetPost(req iwfidl.WorkflowGetRequest, waitIfStillRunning bool) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus) {
	resp, err := s.client.DescribeWorkflowExecution(context.Background(), req.GetWorkflowId(), req.GetWorkflowRunId())
	if err != nil {
		return nil, s.handleError(err)
	}

	var output service.InterpreterWorkflowOutput
	if req.GetNeedsResults() || waitIfStillRunning {
		if resp.Status == service.WorkflowStatusCompleted || waitIfStillRunning {
			err := s.client.GetWorkflowResult(context.Background(), &output, req.GetWorkflowId(), req.GetWorkflowRunId())
			if err != nil {
				return nil, s.handleError(err)
			}
		}
	}

	status := resp.Status
	if waitIfStillRunning {
		// override because when GetWorkflowResult, the workflow is completed
		status = service.WorkflowStatusCompleted
	}

	if err != nil {
		return nil, s.handleError(err)
	}

	return &iwfidl.WorkflowGetResponse{
		WorkflowRunId:  resp.RunId,
		WorkflowStatus: status,
		Results:        output.StateCompletionOutputs,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSearchPost(req iwfidl.WorkflowSearchRequest) (*iwfidl.WorkflowSearchResponse, *ErrorAndStatus) {
	pageSize := int32(1000)
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	resp, err := s.client.ListWorkflow(context.Background(), &ListWorkflowExecutionsRequest{
		PageSize: pageSize,
		Query:    req.GetQuery(),
	})
	if err != nil {
		return nil, s.handleError(err)
	}
	return &iwfidl.WorkflowSearchResponse{
		WorkflowExecutions: resp.Executions,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowResetPost(req iwfidl.WorkflowResetRequest) (*iwfidl.WorkflowResetResponse, *ErrorAndStatus) {
	// TODO https://github.com/indeedeng/iwf/issues/52
	// this is required otherwise Cadence won't accept a reset request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	runId, err := s.client.ResetWorkflow(ctx, req)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &iwfidl.WorkflowResetResponse{
		WorkflowRunId: runId,
	}, nil
}

func (s *serviceImpl) handleError(err error) *ErrorAndStatus {
	// TODO differentiate different error for different codes
	log.Println("encounter error for API", err)
	return &ErrorAndStatus{
		StatusCode: http.StatusInternalServerError,
		Error: iwfidl.ErrorResponse{
			Detail: iwfidl.PtrString(err.Error()),
		},
	}
}
