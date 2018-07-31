package data

import "sort"

// A comparing function will waiting 'true' if a left record less than a right.
type RecordSort func(left, right interface{}) bool

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
