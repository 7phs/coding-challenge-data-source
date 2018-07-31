package commands

import (
	"errors"
	"flag"

	"github.com/7phs/coding-challenge-data-source/places"
)

// Options of the application
type Args struct {
	FileName string  // A path of data file
	Distance float64 // A distance from the office point till an edge of a filter of persons list
}

// Parsing the command-line arguments
func ParseArgs() (args *Args, err error) {
	args = &Args{}

	flag.Float64Var(&args.Distance, "distance", 100., "distance from the "+places.DefaultPlace.String()+", in kilometers")
	flag.Parse()

	if len(flag.Args()) != 1 {
		return nil, errors.New("a data file name: empty")
	} else {
		args.FileName = flag.Args()[0]
	}

	return args, nil
}
