package commands

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-data-source/data"
	"github.com/7phs/coding-challenge-data-source/log"
	"github.com/7phs/coding-challenge-data-source/places"
)

func Root(args *Args) {
	file, err := os.Open(args.FileName)
	if err != nil {
		log.Error("failed to open a data file name '", args.FileName, "': ", err)
		return
	}

	data := data.NewSource(file, places.NewRecordFabric).
		Filter(func(rec interface{}) bool {
			return places.DefaultPlace.Distance(rec.(*places.Record)) <= args.Distance
		}).
		Collect().
		Sort(func(left, right interface{}) bool {
			return left.(*places.Record).Id < right.(*places.Record).Id
		})

	if len(data) == 0 {
		log.Error("no one record found")
		return
	}

	log.Info("Neighbours near", fmt.Sprintf("%.2f", args.Distance), "km(s):")

	for _, rec := range data {
		switch record := rec.(type) {
		case *places.Record:
			log.Infof("#%d, %s", record.Id, record.Name)
		default:
			log.Error("failed to cast a record ", rec, " to expected type")
		}
	}
}
