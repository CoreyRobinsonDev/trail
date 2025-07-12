package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)
var version = "v1.1.0"

// TODO:
// - fix line moving on hitting backspace
// - add (e)xport functionality that writes the trail to a markdown file
// - order -ls output from oldest to newest
// - start command list at 1
// - remove the ability to call -rm alongside other flags
// - color file name in SessionHistory.Load() error
// - (optional) add ability to edit comments



func main() {
	versionPtr := flag.Bool("v", false, "view trail cli version")
	listPtr := flag.Bool("ls", false, "list all past sessions")
	connectPtr := flag.String("conn", "", "connect to a past session")
	exportPtr := flag.String("x", "", "export session to a markdown file")
	removePtr := flag.String("rm", "", "remove a past session")
	flag.Parse()

	if len(*removePtr) != 0 {
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))
		found := false

		if *removePtr == "*" {
			Expect(os.RemoveAll(homedir + "/.config/trail/sessions"))
			Expect(os.MkdirAll(homedir + "/.config/trail/sessions", 0777))
			found = true
		} else {
			for _, file := range files {
				if strings.HasPrefix(file.Name(), *removePtr) {
					Expect(os.Remove(homedir + "/.config/trail/sessions/" + file.Name()))
					found = true
				}
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m a file name being with [%s] could not be found\n", *removePtr)
			os.Exit(1)
		}
	}

	switch(true) {
	case *versionPtr: fmt.Printf("trail %s\n", version)
	case *exportPtr != "":
		var sessionName string
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))

		for _, file := range files {
			if strings.HasPrefix(file.Name(), *exportPtr) {
				sessionName = file.Name()
				break
			}
		}
		sessionHistory := SessionHistory{}
		sessionHistory.Load(sessionName)
		sessionHistory.Export()
	case *listPtr: 
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))
		uniqueIdentifiers := []string{}

		for _, file := range files {
			idx := 1
			fileName := strings.Split(file.Name(), ".")[0]
			for slices.Contains(uniqueIdentifiers, fileName[:idx]) {
				idx++
			}
			uniqueIdentifiers = append(uniqueIdentifiers, fileName[:idx])
			sessionHistory := SessionHistory{}
			Expect(json.Unmarshal(Unwrap(os.ReadFile(homedir + "/.config/trail/sessions/" + file.Name())), &sessionHistory))
			fmt.Printf("\x1b[33m%s\x1b[0m%s", fileName[:idx], fileName[idx:])
			fmt.Printf("\t\x1b[2m%s\x1b[0m\n", sessionHistory.StartTime)
		}
		if len(files) != 0 { 
			fmt.Println("\nrun \x1b[36mtrail -conn \x1b[33m{id}\x1b[0m to continue the session") 
			fmt.Println("run \x1b[36mtrail -rm \x1b[33m{id|\"*\"}\x1b[0m to remove a session") 
			fmt.Println("run \x1b[36mtrail -x \x1b[33m{id|\"*\"}\x1b[0m to export a session") 
		}
	case *connectPtr != "": 
		var sessionName string
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))

		for _, file := range files {
			if strings.HasPrefix(file.Name(), *connectPtr) {
				sessionName = file.Name()
				break
			}
		}
		sessionHistory := SessionHistory{}
		sessionHistory.Load(sessionName)
		sessionHistory.LastModified = time.Now()
		SessionStart(sessionHistory)
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


