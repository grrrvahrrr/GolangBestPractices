package process

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
	RequestBody     string
	FileName        string
	ColumnName      []string
	SearchBody      []string
	SearchParamName []string
	SearchParam     []string
	SearchValue     []string
}

func (r UserRequest) ReadFile() {
	file, err := os.Open(r.FileName)
	if err != nil {
		log.WithError(err).Error("Error openning file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	var indexCol []int
	var indexParam []int
	var counter int

	for counter < 15 {
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

		// log.Info(r.SearchParamName)
		// log.Info(indexParam)

		var sliceToPrint []string
		if indexParam == nil {
			for _, v := range indexCol {
				sliceToPrint = append(sliceToPrint, rec[v])
			}
			log.Info(sliceToPrint)
		} else {
			for i, v := range indexParam {
				switch r.SearchParam[i] {
				case ">":
					if len(indexParam) == 1 {
						recValue, _ := strconv.ParseFloat(rec[v], 64)
						value, _ := strconv.ParseFloat(r.SearchValue[i], 64)
						if recValue > value {
							for _, v := range indexCol {
								sliceToPrint = append(sliceToPrint, rec[v])
							}
							sliceToPrint = append(sliceToPrint, rec[v])
							log.Info(sliceToPrint)
						}
					} else {
						//log.Info("There is more than 1 param")
						var recValueSlice []float64
						for _, v := range indexParam {
							recValue, _ := strconv.ParseFloat(rec[v], 64)
							recValueSlice = append(recValueSlice, recValue)
						}
						//log.Info(recValueSlice)
						var valueSlice []float64
						for i := range r.SearchValue {
							value, _ := strconv.ParseFloat(r.SearchValue[i], 64)
							valueSlice = append(valueSlice, value)
						}
						//log.Info(valueSlice)

						if sliceToPrint == nil {
							for _, v := range indexCol {
								sliceToPrint = append(sliceToPrint, rec[v])
							}
							func() {
								for i := range recValueSlice {
									for range valueSlice {
										//log.Info(recValueSlice)
										//log.Info(valueSlice)
										if recValueSlice[i] < valueSlice[i] {
											return
										}
									}
									if len(sliceToPrint) < (len(indexCol) + len(indexParam)) {
										for _, v := range indexParam {
											sliceToPrint = append(sliceToPrint, rec[v])
										}
									}
								}
								log.Info(sliceToPrint)
							}()

						}
					}

				}

			}

		}

		//counter++
		// if sliceToPrint != nil {
		// 	log.Info(sliceToPrint)
		// }

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

			//log.Info(r.FileName)
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
			//r.SearchParam = r.SearchBody[1]
			r.SearchValue = append(r.SearchValue, r.SearchBody[2])
			for i, v := range r.SearchBody {
				if v == "AND" {
					r.SearchValue = append(r.SearchValue, r.SearchBody[i+3])
				}
			}
			//r.SearchValue = r.SearchBody[2]

			// log.Info(r.SearchBody)
			// log.Info(r.SearchParamName)
			// log.Info(r.SearchParam)
			// log.Info(r.SearchValue)

			break

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

//SELECT SNo, Country/Region FROM /home/deus/Documents/testData/covid_19_data.csv WHERE Confirmed > 100 AND Deaths > 50 AND Recovered > 50
