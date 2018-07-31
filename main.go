package main

import (
	"github.com/7phs/coding-challenge-data-source/commands"
)

const (
	ApplicationName = "Office Neighbors"
	Version         = "0.1"
)

// Variables will initialisation at building time
var (
	GitHash   string // should be uninitialized
	BuildTime string // should be uninitialized
)

// Show an information about the application - a title, a version and building stamps
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
