package places

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/7phs/coding-challenge-data-source/data"
	"github.com/7phs/coding-challenge-data-source/helper"
)

const (
	EarthMidRadius    = 6371.0088
	RadialCoefficient = math.Pi / 180.0
)

var (
	// The default place is a center of an area to find neighbours in it.
	DefaultPlace = (&Record{
		Name:      "Central Office",
		Latitude:  53.339428,
		Longitude: -6.257664,
	}).preCalc()
)

// A record to store an information about a person: id, name and coordinate.
type Record struct {
	Id        int        `json:"user_id"`
	Name      string     `json:"name"`
	Latitude  data.Float `json:"latitude"`
	Longitude data.Float `json:"longitude"`

	// pre-calculated values to prevent repeated calculation
	lat    float64
	long   float64
	cosLat float64
}

// A fabric method to creating a record implementing a validation interface
func NewRecordFabric() data.ValidatedRecord {
	return NewRecord()
}

// A default constructor of a record
func NewRecord() *Record {
	return &Record{}
}

// Pre-calculating a latitude and a longitude in radial coordinates.
func (o *Record) preCalc() *Record {
	o.lat = float64(o.Latitude) * RadialCoefficient
	o.long = float64(o.Longitude) * RadialCoefficient

	o.cosLat = math.Cos(o.lat)

	return o
}

// Stringify a record.
func (o *Record) String() string {
	buf := bytes.NewBufferString("")
	if o.Id > 0 {
		buf.WriteString(fmt.Sprintf("#%d ", o.Id))
	}

	buf.WriteString(fmt.Sprintf("'%s' (%.6f, %.6f)", o.Name, o.Latitude, o.Longitude))

	return buf.String()
}

// Validate fields of a record. Checking empty id, empty name and coordinate are over an edge.
func (o *Record) Validate() error {
	var errList helper.ErrList

	if o.Id <= 0 {
		errList.Add(errors.New("id: empty"))
	}

	if len(o.Name) == 0 {
		errList.Add(errors.New("name: empty"))
	}

	if o.Latitude < -90. || o.Latitude > 90. {
		errList.Add(errors.New("latitude: less than -90 or great than 90"))
	}

	if o.Longitude < -180. || o.Longitude > 180. {
		errList.Add(errors.New("longitude: less than -180 or great than 180"))
	}

	return errList.Finish()
}

// Calculating a great-circle distance between two places.
// Based on a code from http://www.movable-type.co.uk/scripts/latlong.html
func (o *Record) Distance(p *Record) float64 {
	p.preCalc()

	a1 := math.Sin((p.lat - o.lat) / 2)
	a2 := math.Sin((p.long - o.long) / 2)

	a := a1*a1 + a2*a2*o.cosLat*p.cosLat

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthMidRadius * c
}
