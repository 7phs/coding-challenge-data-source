package data

import "sort"

// A comparing function will waiting 'true' if a left record less than a right.
type RecordSort func(left, right interface{}) bool

type RecordMap func(rec interface{}) interface{}

// A helper to collect a record and then sorted it
type RecordList []interface{}

func (o *RecordList) Add(record interface{}) {
	*o = append(*o, record)
}

// Sorting a list of records using a user function to compare records.
func (o RecordList) Sort(sorter RecordSort) RecordList {
	sort.Slice(o, func(i, j int) bool {
		return sorter(o[i], o[j])
	})

	return o
}

// Map each records into another type.
func (o RecordList) Map(mapper RecordMap) RecordList {
	result := make(RecordList, 0, len(o))

	for _, rec := range o {
		result.Add(mapper(rec))
	}

	return result
}
