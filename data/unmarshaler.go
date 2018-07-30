package data

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	supportedUnmarshaling = []RecordUnMarshaling{
		json.Unmarshal,
	}
)

type RecordUnMarshaling func(line []byte, record interface{}) error

func MatchUnMarshaling(line []byte, record interface{}) (RecordUnMarshaling, error) {
	for _, unmarshaling := range supportedUnmarshaling {
		err := unmarshaling(line, record)
		if err == nil {
			return unmarshaling, nil
		}
	}

	return nil, errors.New("unsupported unmarshaling")
}

type Float float64
type AliasFloat Float

func (f *Float) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	data = bytes.Trim(data, "\"")

	return json.Unmarshal(data, (*AliasFloat)(f))
}
