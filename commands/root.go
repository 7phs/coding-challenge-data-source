package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/7phs/coding-challenge-data-source/data"
	"github.com/7phs/coding-challenge-data-source/log"
	"github.com/7phs/coding-challenge-data-source/places"
)

// Loading a list on neighbours which lived around the office using a data source.
// A result list will sort by a person id of neighbours.
func getNeighboursList(reader io.Reader, distance float64) data.RecordList {
	return data.NewSource(reader, places.NewRecordFabric).
		Filter(func(rec interface{}) bool {
			return places.DefaultPlace.Distance(rec.(*places.Record)) <= distance
		}).
		Catch(func(err error) {
			log.Error(err)
		}).
		Collect().
		Sort(func(left, right interface{}) bool {
			return left.(*places.Record).Id < right.(*places.Record).Id
		})
}

// A primary command of the application.
func Root(args *Args) {
	file, err := os.Open(args.FileName)
	if err != nil {
		log.Error("failed to open a data file name:", err)
		return
	}

	neighboursList := getNeighboursList(file, args.Distance)

	if len(neighboursList) == 0 {
		log.Error("No one neighbour found")
		return
	}

	log.Info("Neighbours near", fmt.Sprintf("%.2f", args.Distance), "km(s):")
	for _, rec := range neighboursList {
		switch record := rec.(type) {
		case *places.Record:
			log.Infof("#%d, %s", record.Id, record.Name)
		default:
			log.Error("failed to cast a record ", rec, " to expected type")
		}
	}
}
