package core

import (
	"bytes"
	"encoding/binary"

)

type IOperation interface {
	apply(IStore)
	marshal() ([]byte, error)
	unmarshal([]byte) error
}

type Operation struct {
	name  string
	key   string
	value string
}

func (o Operation) apply(s IStore) interface{} {
	switch o.name {
	case "put":
		return s.ApplyPut(o.key, o.value)
	case "delete":
		return s.ApplyDelete(o.key)
	default:
		panic("unknow command " + o.name)
	}
}

func (o Operation) marshal() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, o)
	return buf.Bytes(), err
}

func (o Operation) unmarshal(b []byte) error {
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, &o)
	return err
}

