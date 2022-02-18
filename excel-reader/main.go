package main

import (
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/xuri/excelize/v2"
)

type requestInput struct {
	Excel   reusable.File `json:"excel"`
	Sheet   string        `json:"sheet"`
	Columns []int         `json:"columns"`
	Start   int           `json:"start"`
}

func networkHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	rf, err := obj.Excel.AsReader()
	if err != nil {
		reusable.ReportError(w, errForCode("file"), err)
	}
	defer rf.Close()

	f, err := excelize.OpenReader(rf)
	if err != nil {
		reusable.ReportError(w, errForCode("file"), err)
	}
	defer f.Close()

	ri.Logger().Infof("open sheet %s", obj.Sheet)

	rows, err := f.GetRows(obj.Sheet)
	if err != nil {
		reusable.ReportError(w, errForCode("file"), err)
		return
	}

	ri.Logger().Infof("starting at row %d from %d rows total", obj.Start, len(rows))

	// test if in range
	start := obj.Start
	if start > len(rows) || len(rows) == 0 {
		reusable.ReportError(w, errForCode("file"), fmt.Errorf("start is larger than number of rows"))
		return
	}

	objs := [][]string{}

	if len(obj.Columns) > 0 {
		ri.Logger().Infof("picking columns %v", obj.Columns)
	} else {
		ri.Logger().Infof("picking all columns")
	}

	for {
		row := rows[start]
		start++

		if len(obj.Columns) > 0 {
			r := []string{}
			for _, a := range obj.Columns {
				r = append(r, row[a])
			}
			objs = append(objs, r)
		} else {
			objs = append(objs, row)
		}

		if start == len(rows) {
			break
		}
	}

	reusable.ReportResult(w, objs)
}

func main() {

	reusable.StartServer(networkHandler, nil)

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.excelreader.%s.error", errCode)
}
