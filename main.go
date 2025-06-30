package main

import (
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

func main() {
	keyEvents := Unwrap(keyboard.GetKeys(10))
	Expect(keyboard.Open())
	defer func() {
		_ = keyboard.Close()
	}()
	sessionHistory := SessionHistory{
		0, 
		Unwrap(os.Stat(historyFile)).ModTime(), 
		[]History{},
	}


	fmt.Println("\x1b[2J\x1b[H\x1b[90mc \x1b[2mto add a comment • \x1b[0m\x1b[90mesc \x1b[2mquit\x1b[0m")
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
								break
							}

						}
					}
				}
			} else {
				if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
					os.Exit(0)
				}
			}
		default:
			if sessionHistory.LastModified.Before(Unwrap(os.Stat(historyFile)).ModTime()) {
				sessionHistory.AddHistory()
				sessionHistory.Render()
			}
		
		}
		time.Sleep(time.Millisecond*200)
	}
}

type History struct {
	Index int
	Cmd string
	Comments []string
}

type SessionHistory struct {
	Index int
	LastModified time.Time
	History []History
}

func (sh *SessionHistory) AddHistory() {
	historyText := strings.Split(string(Unwrap(os.ReadFile(historyFile))), "\n")
	arr := strings.Split(historyText[len(historyText)-2], ";")
	sh.History = append(sh.History, History{sh.Index, arr[len(arr)-1], []string{}})
	sh.Index++
	sh.LastModified = Unwrap(os.Stat(historyFile)).ModTime()
}

func (sh *SessionHistory) AddComment(idx int, content string) {
	if idx < 0 || idx > len(sh.History) {return}
	sh.History[idx].Comments = append(sh.History[idx].Comments, content)
}

func (sh SessionHistory) Render() {
	fmt.Println("\x1b[2J\x1b[H\x1b[90mc \x1b[2mto add a comment • \x1b[0m\x1b[90mesc \x1b[2mquit\x1b[0m")
	fmt.Printf("%s", sh)
}

func (sh SessionHistory) String() string {
	text := ""
	for _, item := range sh.History {
		text += fmt.Sprintf(
			"\x1b[2K\x1b[2m%d\x1b[0m %s\n",
			item.Index,
			item.Cmd,
		)
		for _, Comment := range item.Comments {
			text += fmt.Sprintf("\x1b[2K  • %s\n", Comment)
		}
	}

	return text
}
