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

// WorkflowStateDecideRequest struct for WorkflowStateDecideRequest
type WorkflowStateDecideRequest struct {
	Context Context `json:"context"`
	WorkflowType string `json:"workflowType"`
	WorkflowStateId string `json:"workflowStateId"`
	SearchAttributes []SearchAttribute `json:"searchAttributes,omitempty"`
	QueryAttributes []KeyValue `json:"queryAttributes,omitempty"`
	StateLocalAttributes []KeyValue `json:"stateLocalAttributes,omitempty"`
	CommandResults *CommandResults `json:"commandResults,omitempty"`
}

// NewWorkflowStateDecideRequest instantiates a new WorkflowStateDecideRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowStateDecideRequest(context Context, workflowType string, workflowStateId string) *WorkflowStateDecideRequest {
	this := WorkflowStateDecideRequest{}
	this.Context = context
	this.WorkflowType = workflowType
	this.WorkflowStateId = workflowStateId
	return &this
}

// NewWorkflowStateDecideRequestWithDefaults instantiates a new WorkflowStateDecideRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowStateDecideRequestWithDefaults() *WorkflowStateDecideRequest {
	this := WorkflowStateDecideRequest{}
	return &this
}

// GetContext returns the Context field value
func (o *WorkflowStateDecideRequest) GetContext() Context {
	if o == nil {
		var ret Context
		return ret
	}

	return o.Context
}

// GetContextOk returns a tuple with the Context field value
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetContextOk() (*Context, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Context, true
}

// SetContext sets field value
func (o *WorkflowStateDecideRequest) SetContext(v Context) {
	o.Context = v
}

// GetWorkflowType returns the WorkflowType field value
func (o *WorkflowStateDecideRequest) GetWorkflowType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowType
}

// GetWorkflowTypeOk returns a tuple with the WorkflowType field value
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetWorkflowTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowType, true
}

// SetWorkflowType sets field value
func (o *WorkflowStateDecideRequest) SetWorkflowType(v string) {
	o.WorkflowType = v
}

// GetWorkflowStateId returns the WorkflowStateId field value
func (o *WorkflowStateDecideRequest) GetWorkflowStateId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowStateId
}

// GetWorkflowStateIdOk returns a tuple with the WorkflowStateId field value
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetWorkflowStateIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowStateId, true
}

// SetWorkflowStateId sets field value
func (o *WorkflowStateDecideRequest) SetWorkflowStateId(v string) {
	o.WorkflowStateId = v
}

// GetSearchAttributes returns the SearchAttributes field value if set, zero value otherwise.
func (o *WorkflowStateDecideRequest) GetSearchAttributes() []SearchAttribute {
	if o == nil || o.SearchAttributes == nil {
		var ret []SearchAttribute
		return ret
	}
	return o.SearchAttributes
}

// GetSearchAttributesOk returns a tuple with the SearchAttributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetSearchAttributesOk() ([]SearchAttribute, bool) {
	if o == nil || o.SearchAttributes == nil {
		return nil, false
	}
	return o.SearchAttributes, true
}

// HasSearchAttributes returns a boolean if a field has been set.
func (o *WorkflowStateDecideRequest) HasSearchAttributes() bool {
	if o != nil && o.SearchAttributes != nil {
		return true
	}

	return false
}

// SetSearchAttributes gets a reference to the given []SearchAttribute and assigns it to the SearchAttributes field.
func (o *WorkflowStateDecideRequest) SetSearchAttributes(v []SearchAttribute) {
	o.SearchAttributes = v
}

// GetQueryAttributes returns the QueryAttributes field value if set, zero value otherwise.
func (o *WorkflowStateDecideRequest) GetQueryAttributes() []KeyValue {
	if o == nil || o.QueryAttributes == nil {
		var ret []KeyValue
		return ret
	}
	return o.QueryAttributes
}

// GetQueryAttributesOk returns a tuple with the QueryAttributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetQueryAttributesOk() ([]KeyValue, bool) {
	if o == nil || o.QueryAttributes == nil {
		return nil, false
	}
	return o.QueryAttributes, true
}

// HasQueryAttributes returns a boolean if a field has been set.
func (o *WorkflowStateDecideRequest) HasQueryAttributes() bool {
	if o != nil && o.QueryAttributes != nil {
		return true
	}

	return false
}

// SetQueryAttributes gets a reference to the given []KeyValue and assigns it to the QueryAttributes field.
func (o *WorkflowStateDecideRequest) SetQueryAttributes(v []KeyValue) {
	o.QueryAttributes = v
}

// GetStateLocalAttributes returns the StateLocalAttributes field value if set, zero value otherwise.
func (o *WorkflowStateDecideRequest) GetStateLocalAttributes() []KeyValue {
	if o == nil || o.StateLocalAttributes == nil {
		var ret []KeyValue
		return ret
	}
	return o.StateLocalAttributes
}

// GetStateLocalAttributesOk returns a tuple with the StateLocalAttributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetStateLocalAttributesOk() ([]KeyValue, bool) {
	if o == nil || o.StateLocalAttributes == nil {
		return nil, false
	}
	return o.StateLocalAttributes, true
}

// HasStateLocalAttributes returns a boolean if a field has been set.
func (o *WorkflowStateDecideRequest) HasStateLocalAttributes() bool {
	if o != nil && o.StateLocalAttributes != nil {
		return true
	}

	return false
}

// SetStateLocalAttributes gets a reference to the given []KeyValue and assigns it to the StateLocalAttributes field.
func (o *WorkflowStateDecideRequest) SetStateLocalAttributes(v []KeyValue) {
	o.StateLocalAttributes = v
}

// GetCommandResults returns the CommandResults field value if set, zero value otherwise.
func (o *WorkflowStateDecideRequest) GetCommandResults() CommandResults {
	if o == nil || o.CommandResults == nil {
		var ret CommandResults
		return ret
	}
	return *o.CommandResults
}

// GetCommandResultsOk returns a tuple with the CommandResults field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideRequest) GetCommandResultsOk() (*CommandResults, bool) {
	if o == nil || o.CommandResults == nil {
		return nil, false
	}
	return o.CommandResults, true
}

// HasCommandResults returns a boolean if a field has been set.
func (o *WorkflowStateDecideRequest) HasCommandResults() bool {
	if o != nil && o.CommandResults != nil {
		return true
	}

	return false
}

// SetCommandResults gets a reference to the given CommandResults and assigns it to the CommandResults field.
func (o *WorkflowStateDecideRequest) SetCommandResults(v CommandResults) {
	o.CommandResults = &v
}

func (o WorkflowStateDecideRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["context"] = o.Context
	}
	if true {
		toSerialize["workflowType"] = o.WorkflowType
	}
	if true {
		toSerialize["workflowStateId"] = o.WorkflowStateId
	}
	if o.SearchAttributes != nil {
		toSerialize["searchAttributes"] = o.SearchAttributes
	}
	if o.QueryAttributes != nil {
		toSerialize["queryAttributes"] = o.QueryAttributes
	}
	if o.StateLocalAttributes != nil {
		toSerialize["stateLocalAttributes"] = o.StateLocalAttributes
	}
	if o.CommandResults != nil {
		toSerialize["commandResults"] = o.CommandResults
	}
	return json.Marshal(toSerialize)
}

type NullableWorkflowStateDecideRequest struct {
	value *WorkflowStateDecideRequest
	isSet bool
}

func (v NullableWorkflowStateDecideRequest) Get() *WorkflowStateDecideRequest {
	return v.value
}

func (v *NullableWorkflowStateDecideRequest) Set(val *WorkflowStateDecideRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowStateDecideRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowStateDecideRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowStateDecideRequest(val *WorkflowStateDecideRequest) *NullableWorkflowStateDecideRequest {
	return &NullableWorkflowStateDecideRequest{value: val, isSet: true}
}

func (v NullableWorkflowStateDecideRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowStateDecideRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


