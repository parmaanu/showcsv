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

type autoCompleterType struct {
	historyFile *os.File
	uniqueLines setutils.StringSet

	ReturnAllLines bool
}

func newAutoCompleter() *autoCompleterType {
	absfname, _ := tilde.Expand(gHistoryFile)
	// historyFile = ""
	// if !fileutils.FileExist(absfname) {
	// }
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
				uniqueLines = uniqueLines.AppendIfMissing(line)
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
	if len(line) > 0 {
		ac.uniqueLines = ac.uniqueLines.AppendIfMissing(line)
		ac.historyFile.WriteString(line + "\n")
	}
}

func (ac *autoCompleterType) getFilteredLines(prefix string) []string {
	if len(prefix) == 0 {
		return []string{}
	}
	entries := []string{}
	if ac.ReturnAllLines {
		return ac.uniqueLines[:gShowItemsCount]
	}
	for _, uniqueLine := range ac.uniqueLines {
		if strings.HasPrefix(strings.ToLower(uniqueLine), strings.ToLower(prefix)) {
			entries = append(entries, uniqueLine)
		}
		if len(entries) >= gShowItemsCount {
			return entries
		}
	}
	return entries

	// return algoutils.Filter(ac.uniqueLines, func(uniqueLine string) bool {
	//     return strings.HasPrefix(strings.ToLower(uniqueLine), strings.ToLower(prefix))
	// })
}

// func (ac *autoCompleterType) addAndGetLines(line string) []string {
//     ac.uniqueLines = ac.uniqueLines.AppendIfMissing(line)
//     ac.historyFile.WriteString(line + "\n")
//     return ac.uniqueLines
// }

// func (ac *autoCompleterType) destroy() {
//     ac.historyFile.Close()
// }
