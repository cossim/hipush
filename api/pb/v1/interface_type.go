package v1

import (
	"encoding/json"
	"errors"
	"github.com/gogo/protobuf/proto"
)

var _ proto.Marshaler = (*InterfaceType)(nil)
var _ proto.Unmarshaler = (*InterfaceType)(nil)

func NewInterfaceType(data interface{}) *InterfaceType {
	return &InterfaceType{
		Value: data,
	}
}

type InterfaceType struct {
	Value interface{}
}

func (t InterfaceType) Marshal() ([]byte, error) {
	return json.Marshal(t.Value)
}
func (t *InterfaceType) MarshalTo(data []byte) (n int, err error) {
	return 0, errors.New("not implement")
}
func (t *InterfaceType) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}
func (t *InterfaceType) Size() int {
	return -1
}

// 因为只做JSON的序列化，所以只实现这两个方法就行了
func (t InterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}
func (t *InterfaceType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}
