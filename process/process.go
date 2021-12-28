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
	ReadFile(ctx context.Context) error
}

type ProcessRequest interface {
	ParseRequest(request string, defaultFile string) error
}

type Processer interface {
	ParseRequest(request string, defaultFile string) error
	ReadFile(ctx context.Context) error
}

type Request struct {
	FileName        string
	ColumnName      []string
	SearchBody      []string
	SearchParamName []string
	SearchParam     []string
	SearchValue     []string
}

func (r *Request) ReadFile(ctx context.Context) error {
	file, err := os.Open(r.FileName)
	if err != nil {
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
			}

			var sliceToPrint []string

			if indexParam == nil {
				header = false
				for _, v := range indexCol {
					sliceToPrint = append(sliceToPrint, rec[v])
				}
				fmt.Fprintln(os.Stdout, sliceToPrint)
			} else {
				if !header {
					for _, v := range indexCol {
						sliceToPrint = append(sliceToPrint, rec[v])
					}
					for i, v := range indexParam {
						stringToAdd, err := processSearchParam(r.SearchParam[i], indexParam, rec, v, i, r.SearchValue)
						if err != nil {
							return err
						}
						if stringToAdd != "" {
							sliceToPrint = append(sliceToPrint, stringToAdd)
						}
					}

					if len(sliceToPrint) == (len(indexCol) + len(indexParam)) {
						fmt.Fprintln(os.Stdout, sliceToPrint)
					}
				}
				header = false
			}
		}
	}
}
