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

// checks if the KeyValue type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KeyValue{}

// KeyValue struct for KeyValue
type KeyValue struct {
	Key   *string        `json:"key,omitempty"`
	Value *EncodedObject `json:"value,omitempty"`
}

// NewKeyValue instantiates a new KeyValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeyValue() *KeyValue {
	this := KeyValue{}
	return &this
}

// NewKeyValueWithDefaults instantiates a new KeyValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeyValueWithDefaults() *KeyValue {
	this := KeyValue{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *KeyValue) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeyValue) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *KeyValue) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *KeyValue) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *KeyValue) GetValue() EncodedObject {
	if o == nil || IsNil(o.Value) {
		var ret EncodedObject
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeyValue) GetValueOk() (*EncodedObject, bool) {
	if o == nil || IsNil(o.Value) {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *KeyValue) HasValue() bool {
	if o != nil && !IsNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given EncodedObject and assigns it to the Value field.
func (o *KeyValue) SetValue(v EncodedObject) {
	o.Value = &v
}

func (o KeyValue) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o KeyValue) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Key) {
		toSerialize["key"] = o.Key
	}
	if !IsNil(o.Value) {
		toSerialize["value"] = o.Value
	}
	return toSerialize, nil
}

type NullableKeyValue struct {
	value *KeyValue
	isSet bool
}

func (v NullableKeyValue) Get() *KeyValue {
	return v.value
}

func (v *NullableKeyValue) Set(val *KeyValue) {
	v.value = val
	v.isSet = true
}

func (v NullableKeyValue) IsSet() bool {
	return v.isSet
}

func (v *NullableKeyValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeyValue(val *KeyValue) *NullableKeyValue {
	return &NullableKeyValue{value: val, isSet: true}
}

func (v NullableKeyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeyValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
