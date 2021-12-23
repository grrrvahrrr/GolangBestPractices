package process

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (r *UserRequest) GetRequest() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, Enter an SQL request: SELECT *column_name* FROM *file_name* WHERE *search_parameter* AND *search_parameter*.\nSearch parameters are optional.")

	var err error
	r.RequestBody, err = reader.ReadString('\n')
	if err != nil {
		log.Error(err)
	}
}

func (r *UserRequest) ParseRequest() {
	r.ColumnName = strings.Split(between(r.RequestBody, "SELECT", "FROM"), ",")
	for i := range r.ColumnName {
		r.ColumnName[i] = strings.TrimSpace(r.ColumnName[i])
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

			// r.SearchValue = append(r.SearchValue, r.SearchBody[2])
			// for i, v := range r.SearchBody {
			// 	if v == "AND" {
			// 		r.SearchValue = append(r.SearchValue, r.SearchBody[i+3])
			// 	}
			// }

			//Parse Search value
			if strings.Contains(strings.Join(r.SearchBody, " "), "AND") {

				var paramCounter int = 0
				for _, v := range r.SearchBody {
					if v == "AND" {
						paramCounter++
					}
				}
				for i := 0; i <= paramCounter; i++ {
					var newParam bool
					if strings.Contains(strings.Join(r.SearchBody, " "), "AND") {
						string := strings.TrimSpace(between(strings.Join(r.SearchBody, " "), r.SearchParam[i], "AND"))
						r.SearchValue = append(r.SearchValue, string)
						for j := range r.SearchBody {
							if j < len(r.SearchBody)-1 && !newParam && r.SearchBody[j] == "AND" {
								r.SearchBody = r.SearchBody[j+1 : len(r.SearchBody)]
								newParam = true
							}
						}
					} else {
						string := strings.TrimSpace(after(strings.Join(r.SearchBody, " "), r.SearchParam[i]))
						r.SearchValue = append(r.SearchValue, string)
					}
				}
			} else {
				string := strings.TrimSpace(after(strings.Join(r.SearchBody, " "), r.SearchParam[0]))
				r.SearchValue = append(r.SearchValue, string)
			}
			break

		} else {
			r.FileName = strings.TrimSpace(after(r.RequestBody, "FROM"))
		}
	}
	// log.Info(r.RequestBody)
	// log.Info(r.FileName)
	// log.Info(r.ColumnName)
	// log.Info(len(r.ColumnName))
	// log.Info(r.SearchBody)
	// log.Info(r.SearchParamName)
	// log.Info(r.SearchParam)
	// log.Info(r.SearchValue)

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
