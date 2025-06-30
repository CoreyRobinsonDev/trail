package main

import (
	"fmt"
	"os"
	"runtime"
)

func Expect(err error) {
	_, file, line, ok := runtime.Caller(1)
	if !ok { 
		if err != nil {
			fmt.Fprintf(
				os.Stderr, 
				"%v\n",
				err,
			)
			os.Exit(1)
		}
	} else {
		if err != nil {
			fmt.Fprintf(
				os.Stderr, 
				"\x1b[2m%s\x1b[0m %v\n",
				fmt.Sprintf("(%s:%d): ", file,line),
				err,
			)
			os.Exit(1)
		}
	}
}

func Unwrap[T any](result T, err error) T {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		if err != nil {
			fmt.Fprintf(
				os.Stderr, 
				"%v\n",
				err,
			)
			os.Exit(1)
		}
	} else {
		if err != nil {
			fmt.Fprintf(
				os.Stderr, 
				"\x1b[2m%s\x1b[0m %v\n",
				fmt.Sprintf("(%s:%d): ",file,line),
				err,
			)
			os.Exit(1)
		}
	}
	return result
}
