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



func main() {
	startPtr := flag.Bool("start", false, "start session")
	versionPtr := flag.Bool("v", false, "view trail cli version")
	listPtr := flag.Bool("list", false, "list all past sessions")
	connectPtr := flag.String("connect", "", "connect to a past session")
	flag.Parse()

	switch(true) {
	case *startPtr: 
		sessionHistory := SessionHistory{
			Id: ReverseString(base64.StdEncoding.EncodeToString([]byte(time.Now().String())))[:16],
			StartTime: time.Now(),
			LastModified: Unwrap(os.Stat(historyFile)).ModTime(), 
			History: []History{},
		}
		SessionStart(sessionHistory)
	case *versionPtr: fmt.Println("trail v1.0.0")
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
		if len(files) != 0 { fmt.Println("\nrun \x1b[36mtrail -connect \x1b[33m{id}\x1b[0m to continue the session") }
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
	}
}


