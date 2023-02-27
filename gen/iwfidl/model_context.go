/*
Workflow APIs

This APIs for iwf SDKs to operate workflows

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package iwfidl

import (
	"encoding/json"
)

// checks if the Context type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Context{}

// Context struct for Context
type Context struct {
	WorkflowId               string `json:"workflowId"`
	WorkflowRunId            string `json:"workflowRunId"`
	WorkflowStartedTimestamp int64  `json:"workflowStartedTimestamp"`
	StateExecutionId         string `json:"stateExecutionId"`
	FirstAttemptTimestamp    *int64 `json:"firstAttemptTimestamp,omitempty"`
	Attempt                  *int32 `json:"attempt,omitempty"`
}

// NewContext instantiates a new Context object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewContext(workflowId string, workflowRunId string, workflowStartedTimestamp int64, stateExecutionId string) *Context {
	this := Context{}
	this.WorkflowId = workflowId
	this.WorkflowRunId = workflowRunId
	this.WorkflowStartedTimestamp = workflowStartedTimestamp
	this.StateExecutionId = stateExecutionId
	return &this
}

// NewContextWithDefaults instantiates a new Context object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewContextWithDefaults() *Context {
	this := Context{}
	return &this
}

// GetWorkflowId returns the WorkflowId field value
func (o *Context) GetWorkflowId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowId
}

// GetWorkflowIdOk returns a tuple with the WorkflowId field value
// and a boolean to check if the value has been set.
func (o *Context) GetWorkflowIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowId, true
}

// SetWorkflowId sets field value
func (o *Context) SetWorkflowId(v string) {
	o.WorkflowId = v
}

// GetWorkflowRunId returns the WorkflowRunId field value
func (o *Context) GetWorkflowRunId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowRunId
}

// GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field value
// and a boolean to check if the value has been set.
func (o *Context) GetWorkflowRunIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowRunId, true
}

// SetWorkflowRunId sets field value
func (o *Context) SetWorkflowRunId(v string) {
	o.WorkflowRunId = v
}

// GetWorkflowStartedTimestamp returns the WorkflowStartedTimestamp field value
func (o *Context) GetWorkflowStartedTimestamp() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.WorkflowStartedTimestamp
}

// GetWorkflowStartedTimestampOk returns a tuple with the WorkflowStartedTimestamp field value
// and a boolean to check if the value has been set.
func (o *Context) GetWorkflowStartedTimestampOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowStartedTimestamp, true
}

// SetWorkflowStartedTimestamp sets field value
func (o *Context) SetWorkflowStartedTimestamp(v int64) {
	o.WorkflowStartedTimestamp = v
}

// GetStateExecutionId returns the StateExecutionId field value
func (o *Context) GetStateExecutionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.StateExecutionId
}

// GetStateExecutionIdOk returns a tuple with the StateExecutionId field value
// and a boolean to check if the value has been set.
func (o *Context) GetStateExecutionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StateExecutionId, true
}

// SetStateExecutionId sets field value
func (o *Context) SetStateExecutionId(v string) {
	o.StateExecutionId = v
}

// GetFirstAttemptTimestamp returns the FirstAttemptTimestamp field value if set, zero value otherwise.
func (o *Context) GetFirstAttemptTimestamp() int64 {
	if o == nil || IsNil(o.FirstAttemptTimestamp) {
		var ret int64
		return ret
	}
	return *o.FirstAttemptTimestamp
}

// GetFirstAttemptTimestampOk returns a tuple with the FirstAttemptTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Context) GetFirstAttemptTimestampOk() (*int64, bool) {
	if o == nil || IsNil(o.FirstAttemptTimestamp) {
		return nil, false
	}
	return o.FirstAttemptTimestamp, true
}

// HasFirstAttemptTimestamp returns a boolean if a field has been set.
func (o *Context) HasFirstAttemptTimestamp() bool {
	if o != nil && !IsNil(o.FirstAttemptTimestamp) {
		return true
	}

	return false
}

// SetFirstAttemptTimestamp gets a reference to the given int64 and assigns it to the FirstAttemptTimestamp field.
func (o *Context) SetFirstAttemptTimestamp(v int64) {
	o.FirstAttemptTimestamp = &v
}

// GetAttempt returns the Attempt field value if set, zero value otherwise.
func (o *Context) GetAttempt() int32 {
	if o == nil || IsNil(o.Attempt) {
		var ret int32
		return ret
	}
	return *o.Attempt
}

// GetAttemptOk returns a tuple with the Attempt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Context) GetAttemptOk() (*int32, bool) {
	if o == nil || IsNil(o.Attempt) {
		return nil, false
	}
	return o.Attempt, true
}

// HasAttempt returns a boolean if a field has been set.
func (o *Context) HasAttempt() bool {
	if o != nil && !IsNil(o.Attempt) {
		return true
	}

	return false
}

// SetAttempt gets a reference to the given int32 and assigns it to the Attempt field.
func (o *Context) SetAttempt(v int32) {
	o.Attempt = &v
}

func (o Context) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Context) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["workflowId"] = o.WorkflowId
	toSerialize["workflowRunId"] = o.WorkflowRunId
	toSerialize["workflowStartedTimestamp"] = o.WorkflowStartedTimestamp
	toSerialize["stateExecutionId"] = o.StateExecutionId
	if !IsNil(o.FirstAttemptTimestamp) {
		toSerialize["firstAttemptTimestamp"] = o.FirstAttemptTimestamp
	}
	if !IsNil(o.Attempt) {
		toSerialize["attempt"] = o.Attempt
	}
	return toSerialize, nil
}

type NullableContext struct {
	value *Context
	isSet bool
}

func (v NullableContext) Get() *Context {
	return v.value
}

func (v *NullableContext) Set(val *Context) {
	v.value = val
	v.isSet = true
}

func (v NullableContext) IsSet() bool {
	return v.isSet
}

func (v *NullableContext) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableContext(val *Context) *NullableContext {
	return &NullableContext{value: val, isSet: true}
}

func (v NullableContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableContext) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
