package cache

import (
	"bytes"
	"encoding/gob"
)

type item struct {
	Value interface{}
}

func serializer(value interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	storeValue := item{
		Value: value,
	}
	err := enc.Encode(storeValue)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deserializer(value []byte) (interface{}, error) {
	var res item
	buffer := bytes.NewReader(value)
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}
