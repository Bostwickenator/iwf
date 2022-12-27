package mapper

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
)

func MapToInternalSearchAttributes(attributes []iwfidl.SearchAttribute) (map[string]interface{}, error) {
	if len(attributes) == 0 {
		return nil, nil
	}
	res := map[string]interface{}{}
	for _, attr := range attributes {
		switch attr.GetValueType() {
		case iwfidl.KEYWORD, iwfidl.DATETIME, iwfidl.TEXT:
			res[attr.GetKey()] = attr.GetStringValue()
		case iwfidl.INT:
			res[attr.GetKey()] = attr.GetIntegerValue()
		case iwfidl.BOOL:
			res[attr.GetKey()] = attr.GetBoolValue()
		case iwfidl.DOUBLE:
			res[attr.GetKey()] = attr.GetDoubleValue()
		case iwfidl.KEYWORD_ARRAY:
			res[attr.GetKey()] = attr.GetStringArrayValue()
		default:
			return nil, fmt.Errorf("unsupported search attribute value type %v", attr.GetValueType())
		}
	}
	return res, nil
}

func MapCadenceToIwfSearchAttributes(searchAttributes *shared.SearchAttributes, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (map[string]iwfidl.SearchAttribute, error) {
	if searchAttributes == nil || len(requestedSearchAttributes) == 0 {
		return nil, nil
	}

	result := make(map[string]iwfidl.SearchAttribute, len(requestedSearchAttributes))

	for _, sa := range requestedSearchAttributes {
		key := sa.GetKey()

		field, ok := searchAttributes.IndexedFields[key]
		if !ok {
			continue
		}
		var object interface{}
		err := client.NewValue(field).Get(&object)
		if err != nil {
			return nil, err
		}
		rv, err := mapToIwfSearchAttribute(key, sa.GetValueType(), object)
		if err != nil {
			return nil, err
		}
		result[key] = *rv
	}

	return result, nil
}

func MapTemporalToIwfSearchAttributes(searchAttributes *common.SearchAttributes, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (map[string]iwfidl.SearchAttribute, error) {
	if searchAttributes == nil || len(requestedSearchAttributes) == 0 {
		return nil, nil
	}

	result := make(map[string]iwfidl.SearchAttribute, len(requestedSearchAttributes))

	for _, sa := range requestedSearchAttributes {
		key := sa.GetKey()

		field, ok := searchAttributes.IndexedFields[key]
		if !ok {
			continue
		}
		var object interface{}
		// NOTE: Temporal require search attributes always use default data converter, so we don't need to use the customized one
		err := converter.GetDefaultDataConverter().FromPayload(field, &object)
		if err != nil {
			return nil, err
		}
		rv, err := mapToIwfSearchAttribute(key, sa.GetValueType(), object)
		if err != nil {
			return nil, err
		}
		result[key] = *rv
	}

	return result, nil
}

func mapToIwfSearchAttribute(key string, valueType iwfidl.SearchAttributeValueType, object interface{}) (*iwfidl.SearchAttribute, error) {
	var strVal string
	var intVal int64
	var floatVal float64
	var boolVal bool
	var arrayVal []interface{}
	rv := &iwfidl.SearchAttribute{
		Key:       iwfidl.PtrString(key),
		ValueType: ptr.Any(valueType),
	}
	var ok bool
	var err error
	switch valueType {
	case iwfidl.KEYWORD, iwfidl.DATETIME, iwfidl.TEXT:
		strVal, ok = object.(string)
		rv.StringValue = &strVal
	case iwfidl.INT:
		intVal, ok = object.(int64)
		rv.IntegerValue = &intVal
	case iwfidl.BOOL:
		boolVal, ok = object.(bool)
		rv.BoolValue = &boolVal
	case iwfidl.DOUBLE:
		floatVal, ok = object.(float64)
		rv.DoubleValue = &floatVal
	case iwfidl.KEYWORD_ARRAY:
		arrayVal, ok = object.([]interface{})
		for _, ele := range arrayVal {
			strVal, eleok := ele.(string)
			if !eleok {
				return nil, err
			}
			rv.StringArrayValue = append(rv.StringArrayValue, strVal)
		}
	default:
		return nil, fmt.Errorf("unsupported search attribute value type %v", valueType)
	}
	if !ok {
		return nil, fmt.Errorf("unable to convert value %v to type %v for key %v", object, valueType, key)
	}
	return rv, nil
}