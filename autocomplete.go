package showcsv

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/parmaanu/goutils/setutils"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

const gHistoryFile = "~/.showcsv.cmd.history"
const gShowItemsCount = 10
const gMinLineLength = 3
const gMinPrefixLength = 2

type autoCompleterType struct {
	historyFile *os.File
	uniqueLines setutils.StringSet

	ReturnAllLines bool
}

func newAutoCompleter() *autoCompleterType {
	absfname, _ := tilde.Expand(gHistoryFile)
	f, err := os.OpenFile(absfname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	uniqueLines := setutils.StringSet{}
	data := make([]byte, 1000)
	for {
		data = data[:cap(data)]
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			return nil
		}
		data = data[:n]
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				for _, uniqueLine := range uniqueLines {
					if line == uniqueLine {
						continue
					}
				}
				uniqueLines = append(uniqueLines, line)
			}
		}
	}

	return &autoCompleterType{
		historyFile: f,
		uniqueLines: uniqueLines,
	}
}

func (ac *autoCompleterType) addLine(line string) {
	line = strings.TrimSpace(line)
	if len(line) >= gMinLineLength {
		for _, uniqueLine := range ac.uniqueLines {
			if line == uniqueLine {
				return
			}
		}
		ac.uniqueLines = append(ac.uniqueLines, line)
		ac.historyFile.WriteString(line + "\n")
	}
}

func (ac *autoCompleterType) getFilteredLines(prefix string) []string {
	// TODO, implement a filtering based on command type
	// for example, up key on "/" should show history related to ["/", "?"] command history and similarly for "|" and ":"
	// TODO, show current filename when or directories when tab is pressed and in gSaveFileCmd is active

	filterLines := func(prefix string) []string {
		entries := []string{}
		for idx := len(ac.uniqueLines) - 1; idx >= 0; idx-- {
			uniqueLine := ac.uniqueLines[idx]
			if len(prefix) == 0 || strings.HasPrefix(strings.ToLower(uniqueLine), strings.ToLower(prefix)) {
				entries = append(entries, uniqueLine)
				if len(entries) >= gShowItemsCount {
					break
				}
			}
		}
		return entries
	}

	if ac.ReturnAllLines {
		return filterLines("")
	}

	if len(prefix) < gMinPrefixLength {
		return []string{}
	}

	return filterLines(prefix)
}
