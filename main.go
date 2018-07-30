package main

import (
	"github.com/7phs/coding-challenge-data-source/commands"
)

const (
	ApplicationName = "Dublin Office Neighbors"
	Version         = "0.1"
)

var (
	GitHash   string // should be uninitialized
	BuildTime string // should be uninitialized
)

func ApplicationInfo() string {
	return ApplicationName + " " + Version + " [" + GitHash + "] " + BuildTime
}

func main() {
	args, err := commands.ParseArgs()
	if err != nil {
		commands.Usage(ApplicationInfo())
		return
	}

	commands.Root(args)
}
