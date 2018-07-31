package data

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestMatchUnMarshaling(t *testing.T) {
	testSuites := []*struct {
		in        string
		expected  RecordUnMarshaling
		expectErr bool
	}{
		{in: `{"latitude": "52.986f375", "user_id": 12, "name": "Christina McArdle", "longitude": "-6.043701"}`, expected: json.Unmarshal},
		{in: `<rec user_id="12"><name>Christina McArdle</name><latitude>52.986f375</latitude><longitude>-6.043701</longitude></rec>`, expected: xml.Unmarshal},
		{in: `12, Christina McArdle, 52.986f375, -6.043701`, expectErr: true},
	}

	for _, test := range testSuites {
		exist, err := MatchUnMarshaling([]byte(test.in))
		if err != nil {
			if !test.expectErr {
				t.Error("'"+test.in+"': failed to match an unmarshal:", err)
			}
		} else {
			if test.expectErr {
				t.Error("'" + test.in + "': failed to catch an error")
			} else if fmt.Sprint(exist) != fmt.Sprint(test.expected) {
				t.Error("'"+test.in+"': failed to match unmarshal. Got ", exist, ", but expected is ", test.expected)
			}
		}

	}
}

func TestFloat_UnmarshalJSON(t *testing.T) {
	testSuites := []*struct {
		in        string
		expected  Float
		expectErr bool
	}{
		{in: "", expectErr: true},
		{in: "\"\"", expected: 0.},
		{in: "234", expected: 234.},
		{in: "\"-189734.3234\"", expected: -189734.3234},
		{in: "2f34", expectErr: true},
		{in: "\"2f34\"", expectErr: true},
		{in: "\"text\"", expectErr: true},
	}

	for _, test := range testSuites {
		var exist Float
		err := json.Unmarshal([]byte(test.in), &exist)
		if err != nil {
			if !test.expectErr {
				t.Error("'"+test.in+"': failed to unmarshal:", err)
			}
		} else {
			if test.expectErr {
				t.Error("'" + test.in + "': failed to catch an error")
			} else if exist != test.expected {
				t.Error("'"+test.in+"': failed to unmarshal. Got ", exist, ", but expected is ", test.expected)
			}
		}
	}
}
