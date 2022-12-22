package api

import (
	"context"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
)

type ApiService interface {
	ApiV1WorkflowStartPost(ctx context.Context, request iwfidl.WorkflowStartRequest) (*iwfidl.WorkflowStartResponse, *ErrorAndStatus)
	ApiV1WorkflowSignalPost(ctx context.Context, request iwfidl.WorkflowSignalRequest) *ErrorAndStatus
	ApiV1WorkflowStopPost(ctx context.Context, request iwfidl.WorkflowStopRequest) *ErrorAndStatus
	ApiV1WorkflowGetQueryAttributesPost(ctx context.Context, request iwfidl.WorkflowGetDataObjectsRequest) (*iwfidl.WorkflowGetDataObjectsResponse, *ErrorAndStatus)
	ApiV1WorkflowGetSearchAttributesPost(ctx context.Context, request iwfidl.WorkflowGetSearchAttributesRequest) (*iwfidl.WorkflowGetSearchAttributesResponse, *ErrorAndStatus)
	ApiV1WorkflowGetPost(ctx context.Context, request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus)
	ApiV1WorkflowGetWithWaitPost(ctx context.Context, request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus)
	ApiV1WorkflowSearchPost(ctx context.Context, request iwfidl.WorkflowSearchRequest) (*iwfidl.WorkflowSearchResponse, *ErrorAndStatus)
	ApiV1WorkflowResetPost(ctx context.Context, request iwfidl.WorkflowResetRequest) (*iwfidl.WorkflowResetResponse, *ErrorAndStatus)
	Close()
}

type ErrorAndStatus struct {
	StatusCode int
	Error      iwfidl.ErrorResponse
}

type UnifiedClient interface {
	Close()
	StartInterpreterWorkflow(ctx context.Context, options StartWorkflowOptions, args ...interface{}) (runId string, err error)
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	CancelWorkflow(ctx context.Context, workflowID string, runID string) error
	ListWorkflow(ctx context.Context, request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error)
	QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error
	DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*DescribeWorkflowExecutionResponse, error)
	GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error
	ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error)
}

type StartWorkflowOptions struct {
	ID                       string
	TaskQueue                string
	WorkflowExecutionTimeout time.Duration
	WorkflowIDReusePolicy    *iwfidl.WorkflowIDReusePolicy
	CronSchedule             *string
	RetryPolicy              *iwfidl.RetryPolicy
}

type ListWorkflowExecutionsRequest struct {
	PageSize int32
	Query    string
}

type ListWorkflowExecutionsResponse struct {
	Executions []iwfidl.WorkflowSearchResponseEntry
}

type DescribeWorkflowExecutionResponse struct {
	Status           iwfidl.WorkflowStatus
	RunId            string
	SearchAttributes map[string]iwfidl.SearchAttribute
}
