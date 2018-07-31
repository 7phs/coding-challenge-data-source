package places

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/7phs/coding-challenge-data-source/data"
)

func TestRecord_String(t *testing.T) {
	rec := NewRecord()

	expected := "'' (0.000000, 0.000000)"
	if exist := rec.String(); exist != expected {
		t.Error("failed to stringify an empty record. Got '" + exist + "', but expected is '" + expected + "'")
	}

	rec.Id = 495
	rec.Name = `Test"Name`
	rec.Latitude = 10.10101
	rec.Longitude = -50.50505

	expected = `#495 'Test"Name' (10.101010, -50.505050)`
	if exist := rec.String(); exist != expected {
		t.Error("failed to stringify a record. Got '" + exist + "', but expected is '" + expected + "'")
	}
}

func TestRecord_Validate(t *testing.T) {
	testSuites := []*struct {
		in        data.ValidatedRecord
		errFields []string
	}{
		{
			in:        NewRecordFabric(),
			errFields: []string{"id:", "name:"},
		},
		{
			in:        &Record{Id: 15, Name: "test", Latitude: 10.5, Longitude: -10.5},
			errFields: []string{},
		},
		{
			in:        &Record{Id: 15, Name: "test", Latitude: -100.5, Longitude: 190.5},
			errFields: []string{"latitude:", "longitude:"},
		},
		{
			in:        &Record{Id: 15, Latitude: 100.5, Longitude: -190.5},
			errFields: []string{"name:", "latitude:", "longitude:"},
		},
	}

	for _, test := range testSuites {
		err := test.in.Validate()
		if err == nil && len(test.errFields) > 0 {
			t.Error("failed to catch an error")
			continue
		}

		exist := 0
		for _, field := range test.errFields {
			if strings.Index(err.Error(), field) >= 0 {
				exist++
			}
		}

		if expected := len(test.errFields); exist != expected {
			t.Error("failed to validate a record. Got '" + err.Error() + "', but expected fields '" + strings.Join(test.errFields, ", ") + "'")
		}
	}
}

func TestRecord_Distance(t *testing.T) {
	place1 := &Record{
		Latitude:  52.986375,
		Longitude: -6.043701,
	}
	place1.preCalc()

	expected := &Record{
		lat:    0.92478670244641,
		long:   -0.105482481456074,
		cosLat: 0.60200492254537,
	}
	if fmt.Sprintf("%.6f", place1.lat) != fmt.Sprintf("%.6f", expected.lat) {
		t.Error("failed to pre calculate a latitude. Got ", place1.lat, ", but expected is ", expected.lat)
	}
	if fmt.Sprintf("%.6f", place1.long) != fmt.Sprintf("%.6f", expected.long) {
		t.Error("failed to pre calculate a longitude. Got ", place1.lat, ", but expected is ", expected.long)
	}
	if fmt.Sprintf("%.6f", place1.cosLat) != fmt.Sprintf("%.6f", expected.cosLat) {
		t.Error("failed to pre calculate a cosine of latitude. Got ", place1.cosLat, ", but expected is ", expected.cosLat)
	}

	place2 := &Record{
		Latitude:  53.74452,
		Longitude: -7.11167,
	}

	expectedDistance := 110.125
	existDistance := place1.Distance(place2)
	if math.Abs(existDistance-expectedDistance) > 1e-2 {
		t.Error("failed to calc a distance. Got ", existDistance, ", but expected is ", expectedDistance)
	}

}
