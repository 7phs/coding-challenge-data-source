package commands

import (
	"errors"
	"flag"

	"github.com/7phs/coding-challenge-data-source/places"
)

type Args struct {
	FileName string
	Distance float64
}

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
