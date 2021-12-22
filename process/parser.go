package process

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (r *UserRequest) GetRequest() {
	//SELECT SNo, Country/Region FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 100 AND Deaths < 50 AND Recovered > 50

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, Enter an SQL request: SELECT *column_name* FROM *file_name* WHERE *search_parameter* AND *search_parameter*.\nSearch parameters are optional.")

	var err error
	r.RequestBody, err = reader.ReadString('\n')
	if err != nil {
		log.Error(err)
	}
}

func (r *UserRequest) ParseRequest() {
	r.ColumnName = strings.Fields(between(r.RequestBody, "SELECT", "FROM"))
	for i := range r.ColumnName {
		r.ColumnName[i] = strings.Trim(r.ColumnName[i], ",")
	}
	bodySlice := strings.Fields(r.RequestBody)
	for _, v := range bodySlice {
		if v == "WHERE" {
			r.FileName = strings.TrimSpace(between(r.RequestBody, "FROM", "WHERE"))

			//Parse Search parameters
			r.SearchBody = strings.Fields(after(r.RequestBody, "WHERE"))
			r.SearchParamName = append(r.SearchParamName, r.SearchBody[0])
			for i, v := range r.SearchBody {
				if v == "AND" {
					r.SearchParamName = append(r.SearchParamName, r.SearchBody[i+1])
				}
			}

			r.SearchParam = append(r.SearchParam, r.SearchBody[1])
			for i, v := range r.SearchBody {
				if v == "AND" {
					r.SearchParam = append(r.SearchParam, r.SearchBody[i+2])
				}
			}

			r.SearchValue = append(r.SearchValue, r.SearchBody[2])
			for i, v := range r.SearchBody {
				if v == "AND" {
					r.SearchValue = append(r.SearchValue, r.SearchBody[i+3])
				}
			}

			break

		} else {
			r.FileName = strings.TrimSpace(after(r.RequestBody, "FROM"))
		}
	}

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
