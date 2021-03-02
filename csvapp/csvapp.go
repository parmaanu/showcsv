package main

import (
	"flag"

	"github.com/parmaanu/showcsv"
)

func main() {
	filePtr := flag.String("f", "", "input csv file")
	noHeaderPtr := flag.Bool("n", false, "no header in the input csv file")
	flag.Parse()

	fname := *filePtr
	if len(fname) == 0 {
		flag.PrintDefaults()
		return
	}

	showcsv.DisplayFile(fname, !*noHeaderPtr)
}
