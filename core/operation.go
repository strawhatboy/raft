package core

import (
	"encoding/json"
	"fmt"
)

type IOperation interface {
	apply(IStore)
	marshal() ([]byte, error)
	unmarshal([]byte) error
}

type Operation struct {
	Name  string	`json:"name"`
	Key   string	`json:"key"`
	Value string	`json:"value"`
}

func (o Operation) apply(s IStore) interface{} {
	switch o.Name {
	case "put":
		return s.ApplyPut(o.Key, o.Value)
	case "delete":
		return s.ApplyDelete(o.Key)
	default:
		panic("unknow command " + o.Name)
	}
}

func (o Operation) marshal() ([]byte, error) {
	fmt.Println(fmt.Sprintf("marshaling Operation: %v", o))
	return json.Marshal(o)
}

func (o *Operation) unmarshal(b []byte) error {
	fmt.Println(fmt.Sprintf("Unmarshaling Operation: %v, from bytes: %v", o, b))
	return json.Unmarshal(b, o)
}

