package main

import (
	"flag"
	"fmt"
	"os"
)



func main() {
	startPtr := flag.Bool("start", false, "start session")
					hasIdx := false
					text = ""
	flag.Parse()

	switch(true) {
	case *startPtr: 
		sessionHistory := SessionHistory{
			LastModified: Unwrap(os.Stat(historyFile)).ModTime(), 
			History: []History{},
		}
		SessionStart(sessionHistory)
	case *versionPtr: fmt.Println("trail v1.0.0")
	case *connectPtr != "": fmt.Println("connect")

	}
}

