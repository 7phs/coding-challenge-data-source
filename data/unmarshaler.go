package data

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
)

var (
	// List of supported unmarshallers of a line: json and xml.
	supportedUnmarshaling = []RecordUnMarshaling{
		json.Unmarshal,
		xml.Unmarshal,
	}
)

type RecordUnMarshaling func(line []byte, record interface{}) error

// Trying applying each unmarshaller to catch one which best fit for a line.
func MatchUnMarshaling(line []byte) (RecordUnMarshaling, error) {
	for _, unmarshaling := range supportedUnmarshaling {
		var record interface{}

		err := unmarshaling(line, &record)
		if err == nil {
			return unmarshaling, nil
		}
	}

	return nil, errors.New("unsupported unmarshaling")
}

// A type to unmarshal a string with quotes into a float value.
type Float float64
type AliasFloat Float

func (f *Float) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, (*AliasFloat)(f))
}
