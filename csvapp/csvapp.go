package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

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

	csvFile, err := os.Open(fname)
	if err != nil {
		log.Fatalln("could not open csv file", err)
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("error while reading records from csv", err)
	}

	if len(records) == 0 {
		log.Println("empty csv file")
		return
	}

	header := make([]string, len(records[0]))
	if *noHeaderPtr {
		for i := 0; i < len(records[0]); i++ {
			header[i] = fmt.Sprintf("Col%d", i+1)
		}
		fmt.Println("No header")
	} else {
		header = records[0]
		records = records[1:]
	}

	tableConfig := &showcsv.TableConfig{fname, header, records}
	showcsv.Display(tableConfig)
}
