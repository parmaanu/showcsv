package showcsv_test

import (
	"strconv"
	"testing"

	"github.com/parmaanu/showcsv"
)

func TestShowcsv(t *testing.T) {
	header := []string{"digits", "string", "number"}
	csvData := [][]string{}
	for i := int64(0); i < 100000; i++ {
		csvData = append(csvData, []string{
			strconv.FormatInt(i, 10),
			"A-" + strconv.FormatInt(i+120, 10),
			strconv.FormatInt(i+120, 10),
		})
	}
	tableConfig := &showcsv.TableConfig{"test.csv", header, csvData}
	if err := showcsv.Display(tableConfig); err != nil {
		panic(err)
	}
}
