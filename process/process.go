package process

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ProcessFile interface {
	ReadFile() error
}

type ProcessRequest interface {
	GetReqest() error
	ParseRequest() error
}

type ProcessAll interface {
	ProcessFile
	ProcessRequest
}

type UserRequest struct {
	RequestBody string
	ColumnName  []string
	FileName    string
	SearchName  []string
	SearchParam []string
}

func (r UserRequest) ReadFile() {
	file, err := os.Open(r.FileName)
	if err != nil {
		log.WithError(err).Error("Error openning file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	counter := 0
	var index []int

	//var header bool = true
	for counter < 10 {

		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithError(err).Error("Error reading CSV")
		}
		// do something with read line
		//fmt.Printf("%+v\n", rec)

		for _, iv := range r.ColumnName {
			for j, jv := range rec {
				if jv == iv {
					index = append(index, j)
				}
			}
		}
		//if !header {
		var sliceToPrint []string
		for _, v := range index {
			sliceToPrint = append(sliceToPrint, rec[v])
		}
		log.Info(sliceToPrint)
		counter++
		//}

		//header = false
	}

}

func (r *UserRequest) GetRequest() {
	//SELECT first_name, last_name FROM my.csv WHERE age > 40 AND status = “sick”

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, Enter an SQL request: SELECT *column_name* FROM *file_name* WHERE *search_parameter* AND *search_parameter*.\nSearch parameters are optional.")

	var err error
	r.RequestBody, err = reader.ReadString('\n')
	if err != nil {
		log.Error(err)
	}

	//log.Info(r.RequestBody)
}

func (r *UserRequest) ParseRequest() {
	r.ColumnName = strings.Fields(between(r.RequestBody, "SELECT", "FROM"))
	for i := range r.ColumnName {
		r.ColumnName[i] = strings.Trim(r.ColumnName[i], ",")
		//log.Info(r.ColumnName[i])
	}
	bodySlice := strings.Fields(r.RequestBody)
	for _, v := range bodySlice {
		if v == "WHERE" {
			r.FileName = strings.TrimSpace(between(r.RequestBody, "FROM", "WHERE"))
		} else {
			r.FileName = strings.TrimSpace(after(r.RequestBody, "FROM"))
		}
	}

	//log.Info(r.FileName)
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
