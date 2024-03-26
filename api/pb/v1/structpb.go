package v1

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func ToStructPB(obj interface{}) (*structpb.Struct, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var structValue structpb.Struct
	if err := jsonpb.UnmarshalString(string(jsonBytes), &structValue); err != nil {
		return nil, err
	}

	return &structValue, nil
}

func StructPBToMap(structPB *structpb.Struct) map[string]interface{} {
	result := make(map[string]interface{})

	if structPB == nil {
		return result
	}

	for key, value := range structPB.Fields {
		result[key] = StructValueToInterface(value)
	}

	return result
}

func StructValueToInterface(value *structpb.Value) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.GetKind().(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return v.NumberValue
	case *structpb.Value_StringValue:
		return v.StringValue
	case *structpb.Value_BoolValue:
		return v.BoolValue
	case *structpb.Value_StructValue:
		return StructPBToMap(v.StructValue)
	case *structpb.Value_ListValue:
		var result []interface{}
		for _, val := range v.ListValue.Values {
			result = append(result, StructValueToInterface(val))
		}
		return result
	default:
		return nil
	}
}
