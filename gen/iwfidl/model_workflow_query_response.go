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

// WorkflowQueryResponse struct for WorkflowQueryResponse
type WorkflowQueryResponse struct {
	QueryAttributes []KeyValue `json:"queryAttributes,omitempty"`
}

// NewWorkflowQueryResponse instantiates a new WorkflowQueryResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowQueryResponse() *WorkflowQueryResponse {
	this := WorkflowQueryResponse{}
	return &this
}

// NewWorkflowQueryResponseWithDefaults instantiates a new WorkflowQueryResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowQueryResponseWithDefaults() *WorkflowQueryResponse {
	this := WorkflowQueryResponse{}
	return &this
}

// GetQueryAttributes returns the QueryAttributes field value if set, zero value otherwise.
func (o *WorkflowQueryResponse) GetQueryAttributes() []KeyValue {
	if o == nil || o.QueryAttributes == nil {
		var ret []KeyValue
		return ret
	}
	return o.QueryAttributes
}

// GetQueryAttributesOk returns a tuple with the QueryAttributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowQueryResponse) GetQueryAttributesOk() ([]KeyValue, bool) {
	if o == nil || o.QueryAttributes == nil {
		return nil, false
	}
	return o.QueryAttributes, true
}

// HasQueryAttributes returns a boolean if a field has been set.
func (o *WorkflowQueryResponse) HasQueryAttributes() bool {
	if o != nil && o.QueryAttributes != nil {
		return true
	}

	return false
}

// SetQueryAttributes gets a reference to the given []KeyValue and assigns it to the QueryAttributes field.
func (o *WorkflowQueryResponse) SetQueryAttributes(v []KeyValue) {
	o.QueryAttributes = v
}

func (o WorkflowQueryResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.QueryAttributes != nil {
		toSerialize["queryAttributes"] = o.QueryAttributes
	}
	return json.Marshal(toSerialize)
}

type NullableWorkflowQueryResponse struct {
	value *WorkflowQueryResponse
	isSet bool
}

func (v NullableWorkflowQueryResponse) Get() *WorkflowQueryResponse {
	return v.value
}

func (v *NullableWorkflowQueryResponse) Set(val *WorkflowQueryResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowQueryResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowQueryResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowQueryResponse(val *WorkflowQueryResponse) *NullableWorkflowQueryResponse {
	return &NullableWorkflowQueryResponse{value: val, isSet: true}
}

func (v NullableWorkflowQueryResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowQueryResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


