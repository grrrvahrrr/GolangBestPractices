package process

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type ProcessFile interface {
	ReadFile() error
}

type ProcessRequest interface {
	ParseRequest() error
}

type Request struct {
	FileName        string
	ColumnName      []string
	SearchBody      []string
	SearchParamName []string
	SearchParam     []string
	SearchValue     []string
}

func (r Request) ReadFile(ctx context.Context) error {
	file, err := os.Open(r.FileName)
	if err != nil {
		log.WithError(err).Error("Error openning file")
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	var indexCol []int
	var indexParam []int
	var header bool = true

	for {
		select {
		case <-ctx.Done():
			file.Close()
			log.Info(ctx.Err())
			return nil
		default:
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.WithError(err).Error("Error reading CSV")
				return err
			}
			//Creating indexes for columns and search parameters
			if header {
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
				header = false
			}

			var sliceToPrint []string

			if indexParam == nil {
				for _, v := range indexCol {
					sliceToPrint = append(sliceToPrint, rec[v])
				}
				log.Info(sliceToPrint)
			} else {
				for _, v := range indexCol {
					sliceToPrint = append(sliceToPrint, rec[v])
				}
				for i, v := range indexParam {
					stringToAdd := processSearchParam(r.SearchParam[i], indexParam, rec, v, i, r.SearchValue)
					if stringToAdd != "" {
						sliceToPrint = append(sliceToPrint, stringToAdd)
					}
				}

				if len(sliceToPrint) == (len(indexCol) + len(indexParam)) {
					fmt.Fprintln(os.Stdout, sliceToPrint)
				}
			}
		}
	}
}
