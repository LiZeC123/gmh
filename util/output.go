package util

import (
	"fmt"
	"io"
	"os"
)

func PrintErrorLog(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msg, args...)
	if err != nil {
		panic(err)
	}
}

func PrintToFile(file io.Writer, msg string, args ...interface{}) {
	_, err := fmt.Fprintf(file, msg, args...)
	if err != nil {
		panic(err)
	}
}

func CloseWithLog(o io.Closer) {
	if err := o.Close(); err != nil {
		PrintErrorLog("failed to close: %v\n", err)
	}
}
