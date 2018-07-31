package commands

import (
	"flag"
	"fmt"
)

// Show a usage information of the application.
func Usage(appInfo string) {
	fmt.Println(appInfo)
	fmt.Println("")
	fmt.Println("Usage: office-neighbors <options> data-file")
	flag.PrintDefaults()
	fmt.Println("")
}
