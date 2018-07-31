package commands

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/7phs/coding-challenge-data-source/data"
	"github.com/7phs/coding-challenge-data-source/places"
)

func TestGetNeighbours(t *testing.T) {
	src := []byte(`
		{"latitude": "52.986375", "user_id": 12, "name": "User 1", "longitude": "-6.043701"}
		{"latitude": "51.92893", "user_id": 1, "name": "User 2", "longitude": "-10.27699"}
		{"latitude": "51.8856167", "user_id": 2, "name": "User 3", "longitude": "-10.4240951"}
		{"latitude": "52.3191841", "user_id": 3, "name": "User 4", "longitude": "-8.5072391"}
		{"latitude": "53.807778", "user_id": 28, "name": "User 5", "longitude": "-7.714444"}
		{"latitude": "53.4692815", "user_id": 7, "name": "User 6", "longitude": "-9.436036"}
		{"latitude": "54.0894797", "user_id": 8, "name": "User 7", "longitude": "-6.18671"}
		{"latitude": "53.038056", "user_id": 26, "name": "User 8", "longitude": "-7.653889"}
		{"latitude": "54.1225", "user_id": 27, "name": "User 9", "longitude": "-8.143333"}
		{"latitude": "53.1229599", "user_id": 6, "name": "User 10", "longitude": "-6.2705202"}
	`)

	stringify := func(rec interface{}) interface{} {
		return rec.(fmt.Stringer).String()
	}

	expected := (data.RecordList{
		&places.Record{
			Id: 6, Name: "User 10", Latitude: 53.1229599, Longitude: -6.2705202,
		},
		&places.Record{
			Id: 12, Name: "User 1", Latitude: 52.986375, Longitude: -6.043701,
		},
	}).Map(stringify)

	exist := getNeighboursList(bytes.NewReader(src), 50).Map(stringify)

	if !reflect.DeepEqual(exist, expected) {
		t.Error("failed to get neighbours list. Got ", exist, ", but expected is ", expected)
	}
}
