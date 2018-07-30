package log

import (
	"fmt"
	"os"
)

func Info(v ...interface{}) {
	fmt.Println(v...)
}

func Infof(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func Error(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}
