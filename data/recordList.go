package data

import "sort"

type RecordSort func(left, right interface{}) bool

type RecordList []interface{}

func (o *RecordList) Add(record interface{}) {
	*o = append(*o, record)
}

func (o RecordList) Sort(sorter RecordSort) RecordList {
	sort.Slice(o, func(i, j int) bool {
		return sorter(o[i], o[j])
	})

	return o
}
