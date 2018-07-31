package data

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewSource(t *testing.T) {
	longLine := bytes.NewBufferString("")
	for i := 0; i < 16000; i++ {
		longLine.WriteString("012345678901234567890123456789")
	}
	longLine.WriteString("\n")

	expectedLen := longLine.Len() - 1

	source := NewSource(longLine, nil)
	existLine, err := source.readLine()
	if err != nil {
		t.Error("failed to read line from a buffer: ", err)
	} else if existLen := len(existLine); existLen != expectedLen {
		t.Error("failed to read whole line")
	}
}

type TstRecord struct {
	Id   int    `json:"id" xml:"id,attr"`
	Name string `json:"name" xml:"name"`
}

func (o *TstRecord) Validate() error {
	if o.Id <= 0 {
		return errors.New("id: empty")
	}

	return nil
}

func (o *TstRecord) String() string {
	return fmt.Sprintf("#%d %s", o.Id, o.Name)
}

func TestNewSource_readRecord(t *testing.T) {
	testSuites := []*struct {
		in        []byte
		skip      int
		expected  *TstRecord
		expectErr bool
	}{
		{
			in:       []byte(`{"id":12, "name":"John"}` + "\n" + `{"id":8, "name":"Smith"}`),
			expected: &TstRecord{Id: 12, Name: "John"},
		},
		{
			in:       []byte(`{"id":1, "name":"John"}` + "\n" + `{"id":8, "name":Smith}`),
			skip:     1,   // skip a valid line (that matched unmarshaling)
			expected: nil, // because the second line has wrong format
		},
		{
			in:       []byte(`{"id":0, "name":"John"}` + "\n" + `{"id":8, "name":"Smith"}`),
			expected: nil, // because a record not validated
		},
		{
			in:       []byte(`<record id="8"><name>Smith</name></record>\n<record id="12"><name>John</name</record>`),
			expected: &TstRecord{Id: 8, Name: "Smith"},
		},
		{
			in:        []byte(`8,Smith\n12,John`),
			expectErr: true,
		},
	}

	for i, test := range testSuites {
		source := NewSource(bytes.NewReader(test.in), func() ValidatedRecord {
			return &TstRecord{}
		})

		for i := 0; i < test.skip; i++ {
			source.readRecord()
		}
		exist, err := source.readRecord()
		if err != nil {
			if !test.expectErr {
				t.Error(i, ": failed to read a record: ", err)
			}
		} else {
			if test.expectErr {
				t.Error(i, ": failed to catch an error")
			} else if test.expected == nil {
				if exist != nil {
					t.Error(i, ": failed to read a record. Got ", exist, ", but expected is ", test.expected)
				}
			} else if !reflect.DeepEqual(exist, test.expected) {
				t.Error(i, ": failed to read a record. Got ", exist, ", but expected is ", test.expected)
			}
		}
	}
}

func TestSource_Catch(t *testing.T) {
	data := []byte(`{"id":0, "name":"John"}`)

	catchCount1 := 0
	catchCount2 := 0

	source := NewSource(bytes.NewReader(data), func() ValidatedRecord {
		return &TstRecord{}
	}).Catch(func(err error) {
		if !strings.HasPrefix(err.Error(), "line #") {
			t.Error("an error '" + err.Error() + "' should has a prefix 'line', but hasn't")
		}
		catchCount1++
	}).Catch(func(err error) {
		if !strings.HasPrefix(err.Error(), "line #") {
			t.Error("an error '" + err.Error() + "' should has a prefix 'line', but hasn't")
		}
		catchCount2++
	})

	expectedCount := 4
	for i := 0; i < expectedCount; i++ {
		source.recordIndex = i + 1
		source.raiseError(errors.New("just an error"))
	}
}

func TestSource_Filter(t *testing.T) {
	data := []byte(`{"id":0, "name":"John"}`)

	source := NewSource(bytes.NewReader(data), func() ValidatedRecord {
		return &TstRecord{}
	}).Filter(func(interface{}) bool {
		return false
	}).Filter(func(interface{}) bool {
		return true
	}).Filter(func(interface{}) bool {
		return true
	})

	expected := 3
	if exist := len(source.filters); exist != expected {
		t.Error("failed to add filters. Got ", exist, ", but expected is ", expected)
	}
}

func TestSource_Collect(t *testing.T) {
	data := []byte(`
	{"id":1, "name":"John"}
	{"id":0, "name":"Smith"}
	{"id":13, "name":Ian"}
	{"id":14, "name":"Ivanov"}
	{"id":16, "name":"Bill"}
	{"id":17, "name":"Rodgers"}
	`)

	errCount := 0

	source := NewSource(bytes.NewReader(data), func() ValidatedRecord {
		return &TstRecord{}
	}).Catch(func(error) {
		errCount++
	}).Filter(func(rec interface{}) bool {
		return strings.Index(rec.(*TstRecord).Name, "v") < 0
	}).Filter(func(rec interface{}) bool {
		return strings.Index(rec.(*TstRecord).Name, "d") < 0
	})

	expected := RecordList{
		&TstRecord{Id: 1, Name: "John"},
		&TstRecord{Id: 16, Name: "Bill"},
	}
	exist := source.Collect()
	if !reflect.DeepEqual(exist, expected) {
		t.Error("failed to collect a list. Got ", exist, ", but expected is ", expected)
	}
}
