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

// WorkflowGetSearchAttributesResponse struct for WorkflowGetSearchAttributesResponse
type WorkflowGetSearchAttributesResponse struct {
	SearchAttributes []SearchAttribute `json:"searchAttributes,omitempty"`
}

// NewWorkflowGetSearchAttributesResponse instantiates a new WorkflowGetSearchAttributesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowGetSearchAttributesResponse() *WorkflowGetSearchAttributesResponse {
	this := WorkflowGetSearchAttributesResponse{}
	return &this
}

// NewWorkflowGetSearchAttributesResponseWithDefaults instantiates a new WorkflowGetSearchAttributesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowGetSearchAttributesResponseWithDefaults() *WorkflowGetSearchAttributesResponse {
	this := WorkflowGetSearchAttributesResponse{}
	return &this
}

// GetSearchAttributes returns the SearchAttributes field value if set, zero value otherwise.
func (o *WorkflowGetSearchAttributesResponse) GetSearchAttributes() []SearchAttribute {
	if o == nil || isNil(o.SearchAttributes) {
		var ret []SearchAttribute
		return ret
	}
	return o.SearchAttributes
}

// GetSearchAttributesOk returns a tuple with the SearchAttributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowGetSearchAttributesResponse) GetSearchAttributesOk() ([]SearchAttribute, bool) {
	if o == nil || isNil(o.SearchAttributes) {
		return nil, false
	}
	return o.SearchAttributes, true
}

// HasSearchAttributes returns a boolean if a field has been set.
func (o *WorkflowGetSearchAttributesResponse) HasSearchAttributes() bool {
	if o != nil && !isNil(o.SearchAttributes) {
		return true
	}

	return false
}

// SetSearchAttributes gets a reference to the given []SearchAttribute and assigns it to the SearchAttributes field.
func (o *WorkflowGetSearchAttributesResponse) SetSearchAttributes(v []SearchAttribute) {
	o.SearchAttributes = v
}

func (o WorkflowGetSearchAttributesResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.SearchAttributes) {
		toSerialize["searchAttributes"] = o.SearchAttributes
	}
	return json.Marshal(toSerialize)
}

type NullableWorkflowGetSearchAttributesResponse struct {
	value *WorkflowGetSearchAttributesResponse
	isSet bool
}

func (v NullableWorkflowGetSearchAttributesResponse) Get() *WorkflowGetSearchAttributesResponse {
	return v.value
}

func (v *NullableWorkflowGetSearchAttributesResponse) Set(val *WorkflowGetSearchAttributesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowGetSearchAttributesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowGetSearchAttributesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowGetSearchAttributesResponse(val *WorkflowGetSearchAttributesResponse) *NullableWorkflowGetSearchAttributesResponse {
	return &NullableWorkflowGetSearchAttributesResponse{value: val, isSet: true}
}

func (v NullableWorkflowGetSearchAttributesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowGetSearchAttributesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
