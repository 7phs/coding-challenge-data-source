package commands

import (
	"flag"
	"fmt"
)

func Usage(appInfo string) {
	fmt.Println(appInfo)
	fmt.Println("")
	fmt.Println("Usage: office-neighbors <options> data-file")
	flag.PrintDefaults()
	fmt.Println("")
}
