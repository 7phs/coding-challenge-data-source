package data

import (
	"reflect"
	"strings"
	"testing"
)

func TestRecordList_Add(t *testing.T) {
	var exist RecordList

	var expected RecordList

	if !reflect.DeepEqual(exist, expected) {
		t.Error("failed to match empty list")
	}

	for i := 0; i < 10; i++ {
		exist.Add(i)
	}

	expected = RecordList{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if !reflect.DeepEqual(exist, expected) {
		t.Error("failed to match a list. Got ", exist, ", but expected is ", expected)
	}
}

func TestRecordList_Sort(t *testing.T) {
	type testRec struct {
		id   int
		name string
	}

	exist := (RecordList{
		testRec{id: 0, name: "test"},
		testRec{id: 0, name: "test2"},
		testRec{id: 0, name: "test4"},
		testRec{id: 1, name: "test3"},
		testRec{id: 1, name: "test5"},
		testRec{id: 1, name: "test7"},
	}).Sort(func(left, right interface{}) bool {
		if left.(testRec).id < right.(testRec).id {
			return true
		} else if left.(testRec).id > right.(testRec).id {
			return false
		}

		return strings.Compare(left.(testRec).name, right.(testRec).name) > 0
	})

	expected := RecordList{
		testRec{id: 0, name: "test4"},
		testRec{id: 0, name: "test2"},
		testRec{id: 0, name: "test"},
		testRec{id: 1, name: "test7"},
		testRec{id: 1, name: "test5"},
		testRec{id: 1, name: "test3"},
	}

	if !reflect.DeepEqual(exist, expected) {
		t.Error("failed to sort a list. Got \n", exist, "\n, but expected is \n", expected)
	}
}
