package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/coreyrobinsondev/keyboard"
)

var (
	historyFile = os.Getenv("HISTFILE")
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
					text := ""
					idx := -1
					fmt.Println()

					commentInputLoop: for {
						if text == "" && !hasIdx {
							fmt.Printf("\r\x1b[2KI> \x1b[2menter command index (ex: \x1b[36m0\x1b[0m \x1b[2mls)\x1b[0m\x1b[30D")
						} else if !hasIdx {
							fmt.Printf("\r\x1b[2KI> %s", text)
						} else if text == "" && hasIdx {
							fmt.Printf("\r\x1b[2K%d> \x1b[2menter comment\x1b[0m\x1b[13D", idx)
						} else {
							var lines int
							text = strings.ReplaceAll(text, "\n", "\n    ")
							if lines == 0 {
								fmt.Printf("\r\x1b[2K%d> %s", idx, text)
							} else {
								fmt.Printf("\r\x1b[2K\x1b[%dA%d> %s", lines, idx, text)
							}
							lines = strings.Count(text, "\n")
						}
						c, k, e := keyboard.GetKey()
						if e != nil { fmt.Fprintf(os.Stderr,"error reading input: %v\n", e) }

						switch(k) {
						case keyboard.KeySpace:
							text += " "
						case keyboard.KeyBackspace, keyboard.KeyBackspace2:
							if len(text) == 0 { continue }
							text = text[:len(text)-1]
						case keyboard.KeyCtrlV:
							text += Unwrap(clipboard.ReadAll())
						case keyboard.KeyCtrlC, keyboard.KeyEsc:
							fmt.Print("\r\x1b[2K\x1b[1A")
							break commentInputLoop
						case keyboard.KeyTab:
							text += "    "
						case keyboard.KeyArrowLeft:
							fmt.Print("\x1b[1D")
						case keyboard.KeyArrowRight:
							fmt.Print("\x1b[1C")
						case keyboard.KeyEnter:
							if !hasIdx {
								var err error
								idx, err = strconv.Atoi(text)
								if err == nil {
									hasIdx = true
									text = ""
								} else { idx = -1 }
							} else {
								sessionHistory.AddComment(idx,text)
								sessionHistory.Render()
								sessionHistory.Save()
								break commentInputLoop
							}
						default:
							text += string(c)
						}
					}
				} else if char == 'r' {
					text := ""
					fmt.Println()

					removeInputLoop: for {
						if text == "" {
							fmt.Printf("\r\x1b[2K> \x1b[2menter command index or range of indexes (ex: \x1b[36m0-4\x1b[0m)\x1b[49D")
						} else {
							fmt.Printf("\r\x1b[2K> %s", text)
						}
						c, k, e := keyboard.GetKey()
						if e != nil { fmt.Fprintf(os.Stderr,"error reading input: %v\n", e) }

						switch(k) {
						case keyboard.KeySpace:
							text += " "
						case keyboard.KeyBackspace, keyboard.KeyBackspace2:
							if len(text) == 0 { continue }
							text = text[:len(text)-1]
						case keyboard.KeyCtrlV:
							text += Unwrap(clipboard.ReadAll())
						case keyboard.KeyCtrlC, keyboard.KeyEsc:
							fmt.Print("\r\x1b[2K\x1b[1A")
							break removeInputLoop
						case keyboard.KeyTab:
							text += "    "
						case keyboard.KeyEnter:
							idx, err := strconv.Atoi(text)
							if err != nil {
								if strings.Contains(text, "-") {
									idxes := strings.Split(text, "-")
									if len(idxes) != 2 { 
										fmt.Print("\r\x1b[2K\x1b[1A")
										break removeInputLoop 
									}
									low, e1 := strconv.Atoi(idxes[0])
									high, e2 := strconv.Atoi(idxes[1])
									if e1 != nil || e2 != nil || low > high { 
										fmt.Print("\r\x1b[2K\x1b[1A")
										break removeInputLoop 
									}

									for range high-low+1 {
										sessionHistory.RemoveHistory(low)
									}
								} else { 
									fmt.Print("\r\x1b[2K\x1b[1A")
									break removeInputLoop 
								}
							} else {
								sessionHistory.RemoveHistory(idx)
							}
							sessionHistory.Render()
							sessionHistory.Save()
							break removeInputLoop
						default:
							text += string(c)
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

func (sh *SessionHistory) RemoveHistory(idx int) {
	if idx < 0 || idx >= len(sh.History) { return }
	left := sh.History[:idx]
	right := sh.History[idx+1:]
	sh.History = append(left, right...)
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
	fmt.Println("\x1b[2J\x1b[H\x1b[36mc\x1b[0m \x1b[90mto add a comment • \x1b[36mr\x1b[0m \x1b[90mto remove a command • \x1b[0m\x1b[36mesc\x1b[0m \x1b[90mquit\x1b[0m")
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
		fmt.Fprintf(os.Stderr, "\x1b[2mtrail:\x1b[0m a file name being with [%s] could not be found\n", file)
		os.Exit(1)
	}
	Expect(json.Unmarshal(bytes, &sh))
}

