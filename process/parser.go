package process

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const and string = "AND"

func (r *UserRequest) GetRequest() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please, Enter an SQL request: SELECT *column_name* FROM *file_name* WHERE *search_parameter* AND *search_parameter*.\nSearch parameters are optional.")

	var err error

	r.RequestBody, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	//Add logging for r.Request body to access.log
	//Add check of r.RequestBody for SELECT, FROM return error of incorrect syntax and log it in main
	return nil
}

func (r *UserRequest) ParseRequest() error {
	r.ColumnName = strings.Split(between(r.RequestBody, "SELECT", "FROM"), ",")
	for i := range r.ColumnName {
		r.ColumnName[i] = strings.TrimSpace(r.ColumnName[i])
	}

	if strings.Contains(r.RequestBody, "WHERE") {
		r.FileName = strings.TrimSpace(between(r.RequestBody, "FROM", "WHERE"))
		if r.FileName == "" {
			//Make custom error
			err := fmt.Errorf("no file name")
			return err
		}

		//Parse Search parameters
		r.SearchBody = strings.Fields(after(r.RequestBody, "WHERE"))

		r.SearchParamName = append(r.SearchParamName, r.SearchBody[0])
		for i, v := range r.SearchBody {
			if v == and {
				r.SearchParamName = append(r.SearchParamName, r.SearchBody[i+1])
			}
		}

		r.SearchParam = append(r.SearchParam, r.SearchBody[1])
		for i, v := range r.SearchBody {
			if v == and {
				r.SearchParam = append(r.SearchParam, r.SearchBody[i+2])
			}
		}

		//Parse Single word Seach value
		// r.SearchValue = append(r.SearchValue, r.SearchBody[2])
		// for i, v := range r.SearchBody {
		// 	if v == "AND" {
		// 		r.SearchValue = append(r.SearchValue, r.SearchBody[i+3])
		// 	}
		// }

		//Parse multiword search value
		var err error
		r.SearchValue, err = parseSearchValue(r.SearchBody, r.SearchParam)
		if err != nil {
			return err
		}

	} else {
		r.FileName = strings.TrimSpace(after(r.RequestBody, "FROM"))
		if r.FileName == "" {
			//Make custom error
			err := fmt.Errorf("no file name")
			return err
		}
	}

	return nil
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
	return value[adjustedPos:]
}

func parseSearchValue(searchBody []string, searchParam []string) ([]string, error) {
	var searchValue []string
	if strings.Contains(strings.Join(searchBody, " "), and) {

		var paramCounter int = 0
		for _, v := range searchBody {
			if v == and {
				paramCounter++
			}
		}
		for i := 0; i <= paramCounter; i++ {
			var newParam bool
			if strings.Contains(strings.Join(searchBody, " "), and) {
				string := strings.TrimSpace(between(strings.Join(searchBody, " "), searchParam[i], and))
				searchValue = append(searchValue, string)
				for j := range searchBody {
					if j < len(searchBody)-1 && !newParam && searchBody[j] == and {
						searchBody = searchBody[j+1:]
						newParam = true
					}
				}
			} else {
				string := strings.TrimSpace(after(strings.Join(searchBody, " "), searchParam[i]))
				searchValue = append(searchValue, string)
			}
		}
	} else {
		string := strings.TrimSpace(after(strings.Join(searchBody, " "), searchParam[0]))
		searchValue = append(searchValue, string)
	}
	if searchValue == nil {
		//Make custom error
		err := fmt.Errorf("no search values")
		return nil, err
	}
	return searchValue, nil
}
