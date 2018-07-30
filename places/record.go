package places

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/7phs/coding-challenge-data-source/data"
)

const (
	EarthMidRadius    = 6371.0088
	RadialCoefficient = math.Pi / 180.0
)

var (
	DefaultPlace = (&Record{
		Name:      "Intercom Dublin Office",
		Latitude:  53.339428,
		Longitude: -6.257664,
	}).preCalc()
)

type Record struct {
	Id        int        `json:"user_id"`
	Name      string     `json:"name"`
	Latitude  data.Float `json:"latitude"`
	Longitude data.Float `json:"longitude"`

	lat    float64
	long   float64
	cosLat float64
}

func NewRecordFabric() data.ValidatedRecord {
	return NewRecord()
}

func NewRecord() *Record {
	return &Record{}
}

func (o *Record) preCalc() *Record {
	o.lat = float64(o.Latitude) * RadialCoefficient
	o.long = float64(o.Longitude) * RadialCoefficient

	o.cosLat = math.Cos(o.lat)

	return o
}

func (o *Record) String() string {
	buf := bytes.NewBufferString("")
	if o.Id > 0 {
		buf.WriteString(fmt.Sprintf("#%d ", o.Id))
	}

	buf.WriteString(fmt.Sprintf("'%s' (%.6f, %.6f)", o.Name, o.Latitude, o.Longitude))

	return buf.String()
}

func (o *Record) Validate() error {
	if o.Id <= 0 {
		return errors.New("id: empty")
	}

	if len(o.Name) == 0 {
		return errors.New("name: empty")
	}

	if o.Latitude < -90. || o.Latitude > 90. {
		return errors.New("latitude: less than -90 or great than 90")
	}

	if o.Longitude < -180. || o.Longitude > 180. {
		return errors.New("longitude: less than -180 or great than 180")
	}

	return nil
}

// http://www.movable-type.co.uk/scripts/latlong.html
func (o *Record) Distance(p *Record) float64 {
	p.preCalc()

	a1 := math.Sin((p.lat - o.lat) / 2)
	a2 := math.Sin((p.long - o.long) / 2)

	a := a1*a1 + a2*a2*o.cosLat*p.cosLat

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthMidRadius * c
}
