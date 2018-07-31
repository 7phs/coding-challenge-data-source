package log

// A package as a wrapper of logging functions.

import (
	"fmt"
	"os"
)

// Dump a log message to a stdout
func Info(v ...interface{}) {
	fmt.Println(v...)
}

// Dump a formatted log message to a stdout
func Infof(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

// Dump a log message to a stderr
func Error(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}
