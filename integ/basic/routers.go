/*
 * Workflow APIs
 *
 * This APIs for iwf SDKs to operate workflows
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package basic

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	WorkflowType = "basic"
	State1       = "S1"
	State2       = "S2"
)

// NewBasicWorkflow returns a new gin server.
func NewBasicWorkflow() *gin.Engine {
	router := gin.Default()

	handler := newHandler()

	router.POST(service.StateStartApi, handler.apiV1WorkflowStateStart)
	router.POST(service.StateDecideApi, handler.apiV1WorkflowStateDecide)

	return router
}

type handler struct{}

func newHandler() *handler {
	return &handler{}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) apiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		// basic workflow go straight to decide methods without any commands
		if req.GetWorkflowStateId() == State1 || req.GetWorkflowStateId() == State2 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.PtrString(service.DeciderTypeAllCommandCompleted),
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) apiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		if req.GetWorkflowStateId() == State1 {
			// go to S2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: iwfidl.PtrString(State2),
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			// go to complete
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: iwfidl.PtrString(service.CompletingWorkflowStateId),
						},
					},
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}
