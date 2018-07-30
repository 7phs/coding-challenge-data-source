package commands

import (
	"flag"
	"fmt"
)

func Usage(appInfo string) {
	fmt.Println(appInfo)
	fmt.Println("Usage: dublin-office-neighbors <options> data-file")
	flag.PrintDefaults()
}
