package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

var (
	historyFile = os.Getenv("HISTFILE")
	text string
	idx = -1
)

func SessionStart(sessionHistory SessionHistory) {
	homedir := Unwrap(os.UserHomeDir())
	Expect(os.MkdirAll(homedir + "/.config/trail/sessions", 0777))
	keyEvents := Unwrap(keyboard.GetKeys(10))
	Expect(keyboard.Open())
	defer keyboard.Close()


	sessionHistory.Render()
	for {
		select {
		case event := <-keyEvents:
			char := event.Rune
			key := event.Key

			if char != 0 {
				if char == 'c' {
					hasIdx := false
					text = ""
					idx = -1
					fmt.Println()
					for {
						if text == "" && !hasIdx {
							fmt.Printf("\r\x1b[2KI> \x1b[2menter command index (ex: \x1b[36m0\x1b[0m \x1b[2mls)\x1b[0m\x1b[30D")
						} else if !hasIdx {
							fmt.Printf("\r\x1b[2KI> %s", text)
						} else if text == "" && hasIdx {
							fmt.Printf("\r\x1b[2K%d> \x1b[2menter comment\x1b[0m\x1b[13D", idx)
						} else {
							fmt.Printf("\r\x1b[2K%d> %s", idx, text)
						}
						c, k, e := keyboard.GetKey()
						if e != nil { fmt.Fprintf(os.Stderr,"error reading input: %v", e) }
						if c != 0 {
							text += string(c)
						} else if k == keyboard.KeySpace {
							text += " "
						} else if k == keyboard.KeyBackspace2 || k == keyboard.KeyBackspace {
							if len(text) == 0 { continue }
							text = text[:len(text)-1]
						}
						if k == keyboard.KeyCtrlC || k == keyboard.KeyEsc {
							fmt.Print("\r\x1b[2K\x1b[1A")
							break
						} else if k == keyboard.KeyEnter {
							if !hasIdx {
								var err error
								idx, err = strconv.Atoi(text)
								if err == nil {
									hasIdx = true
									text = ""
								} else { idx = -1 }
							} else {
								sessionHistory.AddComment(idx,text)
								fmt.Print("\r\x1b[2K\x1b[2A")
								sessionHistory.Render()
								sessionHistory.Save()
								break
							}
						}
					}
				}
			} else {
				if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
					sessionHistory.Save()
					os.Exit(0)
				}
			}
		default:
			if sessionHistory.LastModified.Before(Unwrap(os.Stat(historyFile)).ModTime()) {
				sessionHistory.AddHistory()
				sessionHistory.Render()
				sessionHistory.Save()
			}
		
		}
		time.Sleep(time.Millisecond*100)
	}
}

type History struct {
	Cmd string `json:"cmd"`
	Comments []string `json:"comments"`
}

type SessionHistory struct {
	Id string `json:"id"`
	StartTime time.Time `json:"startTime"`
	LastModified time.Time `json:"lastModified"`
	History []History `json:"history"`
}

func (sh *SessionHistory) AddHistory() {
	historyText := strings.Split(string(Unwrap(os.ReadFile(historyFile))), "\n")
	arr := strings.Split(historyText[len(historyText)-2], ";")
	sh.History = append(sh.History, History{arr[len(arr)-1], []string{}})
	sh.LastModified = Unwrap(os.Stat(historyFile)).ModTime()
}

func (sh *SessionHistory) AddComment(idx int, content string) {
	if idx < 0 || idx > len(sh.History) {return}
	sh.History[idx].Comments = append(sh.History[idx].Comments, content)
}

func (sh SessionHistory) Render() {
	fmt.Println("\x1b[2J\x1b[H\x1b[36mc\x1b[0m \x1b[90mto add a comment • \x1b[0m\x1b[36mesc\x1b[0m \x1b[90mquit\x1b[0m")
	for i, item := range sh.History {
		fmt.Printf(
			"\x1b[2K\x1b[2m%d\x1b[0m %s\n",
			i,
			item.Cmd,
		)
		for _, comment := range item.Comments {
			fmt.Printf("\x1b[2K  • %s\n", comment)
		}
	}
}


func (sh SessionHistory) Save() {
	if len(sh.History) == 0 {return}
	homedir := Unwrap(os.UserHomeDir())
	Expect(os.WriteFile(
		fmt.Sprintf("%s/.config/trail/sessions/%s.json", homedir, sh.Id),
		Unwrap(json.MarshalIndent(sh, "", "\t")),
		0666,
	))
}

func (sh *SessionHistory) Load(file string) {
	homedir := Unwrap(os.UserHomeDir())
	bytes, err := os.ReadFile(homedir + "/.config/trail/sessions/" + file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m %s could not be found\n", file)
		os.Exit(1)
	}
	Expect(json.Unmarshal(bytes, &sh))
}

