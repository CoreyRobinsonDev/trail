package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)
var version = "v1.2.0"

// TODO:
// - fix line moving on hitting backspace
// - (optional) add ability to edit comments



func main() {
	versionPtr := flag.Bool("v", false, "view trail cli version")
	listPtr := flag.Bool("ls", false, "list all past sessions")
	connectPtr := flag.String("conn", "", "connect to a past session")
	exportPtr := flag.String("x", "", "export session to a markdown file")
	removePtr := flag.String("rm", "", "remove a past session")
	flag.Parse()


	switch(true) {
	case *versionPtr: fmt.Printf("trail %s\n", version)
	case *removePtr != "":
		homedir := Unwrap(os.UserHomeDir())

		if *removePtr == "*" {
			Expect(os.RemoveAll(homedir + "/.config/trail/sessions"))
			Expect(os.MkdirAll(homedir + "/.config/trail/sessions", 0777))
			return
		} else {
			sessions := GetSessions()

			for _, session := range sessions {
				if strings.HasPrefix(session.Id, *removePtr) {
					Expect(os.Remove(homedir + "/.config/trail/sessions/" + session.Id + ".json"))
					return
				}
			}
		}

		fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m a file name being with [%s] could not be found\n", *removePtr)
		os.Exit(1)
	case *exportPtr != "":
		sessions := GetSessions()

		for _, session := range sessions {
			if strings.HasPrefix(session.Id, *exportPtr) {
				session.Export()
				return
			}
		}

		fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m a file name beginning with \x1b[36m%s\x1b[0m could not be found\n", *exportPtr)
		os.Exit(1)

	case *listPtr: 
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))
		uniqueIdentifiers := []string{}
		sessions := GetSessions()

		for _, session := range sessions {
			idx := 1
			for slices.Contains(uniqueIdentifiers, session.Id[:idx]) {
				idx++
			}
			uniqueIdentifiers = append(uniqueIdentifiers, session.Id[:idx])
			fmt.Printf("\x1b[33m%s\x1b[0m%s", session.Id[:idx], session.Id[idx:])
			fmt.Printf("\t\x1b[2m%s\x1b[0m\n", session.StartTime)
		}

		if len(files) != 0 { 
			fmt.Println("\nrun \x1b[36mtrail -conn \x1b[33m{id}\x1b[0m to continue the session") 
			fmt.Println("run \x1b[36mtrail -rm \x1b[33m{id|\"*\"}\x1b[0m to remove a session") 
			fmt.Println("run \x1b[36mtrail -x \x1b[33m{id}\x1b[0m to export a session") 
		}
	case *connectPtr != "": 
		sessions := GetSessions()

		for _, session := range sessions {
			if strings.HasPrefix(session.Id, *connectPtr) {
				session.LastModified = time.Now()
				SessionStart(session)
				return
			}
		}

		fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m a file name beginning with \x1b[36m%s\x1b[0m could not be found\n", *connectPtr)
		os.Exit(1)
	default:
		sessionHistory := SessionHistory{
			Id: ReverseString(base64.StdEncoding.EncodeToString([]byte(time.Now().String())))[:16],
			StartTime: time.Now(),
			LastModified: Unwrap(os.Stat(historyFile)).ModTime(), 
			Commands: []Command{},
		}
		SessionStart(sessionHistory)
	}
}


