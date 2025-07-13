package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
					HandleComment(sessionHistory)
				} else if char == 'r' {
					HandleRemove(sessionHistory)
				} else if char == 't' {
					HandleTitleChange(sessionHistory)
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


func HandleComment(sessionHistory SessionHistory) {
	hasIdx := false
	cursorX := 0
	text := ""
	idx := -1
	fmt.Println()

	for {
		// print input line
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

		// move cursor
		if cursorX < 0 {
			fmt.Printf("\x1b[%dD", cursorX*-1)
		}

		// blocking
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
			return
		case keyboard.KeyTab:
			text += "    "
		case keyboard.KeyArrowLeft:
			if cursorX*-1 >= len(text) { continue }
			cursorX--
		case keyboard.KeyArrowRight:
			if cursorX >= 0 { continue }
			cursorX++
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
				return
			}
		default:
			text += string(c)
		}
	}
}

func HandleRemove(sessionHistory SessionHistory) {
	cursorX := 0
	text := ""
	fmt.Println()

	for {
		if text == "" {
			fmt.Printf("\r\x1b[2K> \x1b[2menter command index or range of indexes (ex: \x1b[36m0-4\x1b[0m)\x1b[49D")
		} else {
			fmt.Printf("\r\x1b[2K> %s", text)
		}

		// move cursor
		if cursorX < 0 {
			fmt.Printf("\x1b[%dD", cursorX*-1)
		}

		// blocking
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
			return
		case keyboard.KeyArrowLeft:
			if cursorX*-1 >= len(text) { continue }
			cursorX--
		case keyboard.KeyArrowRight:
			if cursorX >= 0 { continue }
			cursorX++
		case keyboard.KeyTab:
			text += "    "
		case keyboard.KeyEnter:
			idx, err := strconv.Atoi(text)
			if err != nil {
				if strings.Contains(text, "-") {
					idxes := strings.Split(text, "-")
					if len(idxes) != 2 { 
						fmt.Print("\r\x1b[2K\x1b[1A")
						return
					}
					low, e1 := strconv.Atoi(idxes[0])
					high, e2 := strconv.Atoi(idxes[1])
					if e1 != nil || e2 != nil || low > high { 
						fmt.Print("\r\x1b[2K\x1b[1A")
						return
					}

					for range high-low+1 {
						sessionHistory.RemoveHistory(low)
					}
				} else { 
					fmt.Print("\r\x1b[2K\x1b[1A")
					return
				}
			} else {
				sessionHistory.RemoveHistory(idx)
			}
			sessionHistory.Render()
			sessionHistory.Save()
			return
		default:
			text += string(c)
		}
	}
}

func HandleTitleChange(sessionHistory SessionHistory) {
	text := ""
	cursorX := 0
	fmt.Println()

	loop: for {
		if text == "" {
			fmt.Printf("\r\x1b[2K> \x1b[2menter title\x1b[0m\x1b[11D")
		} else {
			fmt.Printf("\r\x1b[2K> %s", text)
		}

		// move cursor
		if cursorX < 0 {
			fmt.Printf("\x1b[%dD", cursorX*-1)
		}

		// blocking
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
			return
		case keyboard.KeyArrowLeft:
			if cursorX*-1 >= len(text) { continue }
			cursorX--
		case keyboard.KeyArrowRight:
			if cursorX >= 0 { continue }
			cursorX++
		case keyboard.KeyTab:
			text += "    "
		case keyboard.KeyEnter:
			sessions := GetSessions()
			for _, session := range sessions {
				if session.Id == text {
					fmt.Printf("\r\x1b[2K> %s", text)
					continue loop
				}
			}
			homedir := Unwrap(os.UserHomeDir())
			// I don't care if this errors lmao
			os.Remove(homedir + "/.config/trail/sessions/" + sessionHistory.Id + ".json")

			sessionHistory.Id = text
			sessionHistory.Render()
			sessionHistory.Save()
			return
		default:
			text += string(c)
		}
	}
}

func GetSessions() []SessionHistory {
		homedir := Unwrap(os.UserHomeDir())
		files := Unwrap(os.ReadDir(homedir + "/.config/trail/sessions"))
		sessions := []SessionHistory{}
		
		for _, file := range files {
			sessionHistory := SessionHistory{}
			Expect(json.Unmarshal(Unwrap(os.ReadFile(homedir + "/.config/trail/sessions/" + file.Name())), &sessionHistory))
			sessions = append(sessions, sessionHistory)
		}

		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].StartTime.Before(sessions[j].StartTime)
		})

	return sessions
}

type Command struct {
	Name string `json:"cmd"`
	Comments []string `json:"comments"`
}

type SessionHistory struct {
	Id string `json:"id"`
	StartTime time.Time `json:"startTime"`
	LastModified time.Time `json:"lastModified"`
	Commands []Command `json:"history"`
}

func (sh *SessionHistory) RemoveHistory(idx int) {
	if idx < 0 || idx >= len(sh.Commands) { return }
	left := sh.Commands[:idx]
	right := sh.Commands[idx+1:]
	sh.Commands = append(left, right...)
}

func (sh *SessionHistory) AddHistory() {
	historyText := strings.Split(string(Unwrap(os.ReadFile(historyFile))), "\n")
	arr := strings.Split(historyText[len(historyText)-2], ";")
	sh.Commands = append(sh.Commands, Command{arr[len(arr)-1], []string{}})
	sh.LastModified = Unwrap(os.Stat(historyFile)).ModTime()
}

func (sh *SessionHistory) AddComment(idx int, content string) {
	if idx < 0 || idx > len(sh.Commands) {return}
	sh.Commands[idx].Comments = append(sh.Commands[idx].Comments, content)
}

func (sh SessionHistory) Render() {
	fmt.Println("\x1b[2J\x1b[H\x1b[36mc\x1b[0m \x1b[90mto add a comment • \x1b[36mr\x1b[0m \x1b[90mto remove a command • \x1b[36mt\x1b[0m \x1b[90mto set a session title • \x1b[0m\x1b[36mesc\x1b[0m \x1b[90mquit\x1b[0m")
	for i, item := range sh.Commands {
		fmt.Printf(
			"\x1b[2K\x1b[2m%d\x1b[0m %s\n",
			i+1,
			item.Name,
		)
		for _, comment := range item.Comments {
			fmt.Printf("\x1b[2K  • %s\n", comment)
		}
	}
}


func (sh SessionHistory) Save() {
	if len(sh.Commands) == 0 {return}
	homedir := Unwrap(os.UserHomeDir())
	Expect(os.WriteFile(
		fmt.Sprintf("%s/.config/trail/sessions/%s.json", homedir, sh.Id),
		Unwrap(json.MarshalIndent(sh, "", "\t")),
		0666,
	))
}

func (sh *SessionHistory) Load(file string) {
	homedir := Unwrap(os.UserHomeDir())
	bytes := Unwrap(os.ReadFile(homedir + "/.config/trail/sessions/" + file))
	Expect(json.Unmarshal(bytes, &sh))
}

func (sh SessionHistory) Export() {
	file := Unwrap(os.Create("./" + sh.Id + ".md"))
	defer file.Close()

	file.WriteString(fmt.Sprintf(`---
title: %s
date: %v
---
`, sh.Id, sh.StartTime))

	for _, command := range sh.Commands {
		file.WriteString(fmt.Sprintf("1. %s\n", command.Name))

		for _, comment := range command.Comments {
			file.WriteString(fmt.Sprintf("\t- %s\n", comment))
		}
	}
}
