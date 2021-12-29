package process

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func (r *MockedRequest) ReadFile(ctx context.Context, w io.Writer) error {
	var indexCol []int
	var indexParam []int
	var eof bool

	for !eof {
		select {
		case <-ctx.Done():
			log.Info(ctx.Err())
			return nil
		default:
			eof = true
			//Creating indexes for columns and search parameters

			rec := []string{"SNo", "ObservationDate", "Province/State", `Country/Region`, "Last Update", "Confirmed", "Deaths", "Recovered"}
			for _, iv := range r.ColumnName {
				for j, jv := range rec {
					if jv == iv {
						indexCol = append(indexCol, j)
					}
				}
			}
			for _, iv := range r.SearchParamName {
				for j, jv := range rec {
					if jv == iv {
						indexParam = append(indexParam, j)
					}
				}
			}

			var sliceToPrint []string
			rec = []string{"306427", "date", "Zhejiang", `Mainland China`, "update", "1364.0", "1.0", "1324.0"}
			if indexParam == nil {
				for _, v := range indexCol {
					sliceToPrint = append(sliceToPrint, rec[v])
				}
				fmt.Fprintln(os.Stdout, sliceToPrint)
			} else {

				for _, v := range indexCol {
					sliceToPrint = append(sliceToPrint, rec[v])
				}
				for i := range indexParam {
					stringToAdd, err := processSearchParam(r.SearchParam[i], indexParam, rec, i, r.SearchValue)
					if err != nil {
						return err
					}
					if stringToAdd != "" {
						sliceToPrint = append(sliceToPrint, stringToAdd)
					}
				}

				if len(sliceToPrint) == (len(indexCol) + len(indexParam)) {
					fmt.Fprintln(w, sliceToPrint)
				}

			}
		}
	}
	return nil
}

func TestMockedReadFile(t *testing.T) {
	var p ProcessFile = &MockedRequest{
		FileName:        "testdir/test.csv",
		ColumnName:      []string{"SNo", "Province/State"},
		SearchBody:      []string{"Country/Region", "=", `"Mainland`, `China"`, "AND", "Confirmed", ">", "100", "AND", "Deaths", "<", "50", "AND", "Recovered", ">", "20"},
		SearchParamName: []string{"Country/Region", "Confirmed", "Deaths", "Recovered"},
		SearchParam:     []string{"=", ">", "<", ">"},
		SearchValue:     []string{`"Mainland China"`, "100", "50", "20"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var buf bytes.Buffer

	err := p.ReadFile(ctx, &buf)
	if err != nil {
		t.Error(err)
	}

	expectedOutput := "[306427 Zhejiang Mainland China 1364.0 1.0 1324.0]\n"
	output := buf.String()
	if !assert.Equal(t, expectedOutput, output) {
		t.Error("Output was different")
	}

}
