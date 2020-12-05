package showcsv_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/parmaanu/showcsv"
)

var renderTui = false

func init() {
	_, renderTui = os.LookupEnv("SHOWCSV_RENDER_TUI")
}

func TestShowcsv(t *testing.T) {
	if !renderTui {
		return
	}
	header := []string{"digits", "string", "number"}
	csvData := [][]string{}
	for i := int64(0); i < 100000; i++ {
		csvData = append(csvData, []string{
			strconv.FormatInt(i, 10),
			"A-" + strconv.FormatInt(i+120, 10),
			strconv.FormatInt(i+120, 10),
		})
	}
	tableConfig := &showcsv.TableConfig{"example.csv", header, csvData}
	if err := showcsv.Display(tableConfig); err != nil {
		panic(err)
	}
}
