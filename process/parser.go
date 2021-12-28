package process

import (
	"fmt"
	"strings"
)

const and string = "AND"
const from string = "FROM"
const where string = "WHERE"
const sel string = "SELECT"

func (r *Request) ParseRequest(request string, defaultFile string) error {
	const def string = "default"
	r.ColumnName = strings.Split(between(request, sel, from), ",")
	for i := range r.ColumnName {
		r.ColumnName[i] = strings.TrimSpace(r.ColumnName[i])
	}

	if strings.Contains(request, where) {
		r.FileName = strings.TrimSpace(between(request, from, where))
		if r.FileName == "" {
			err := fmt.Errorf("no file name")
			return err
		}
		if r.FileName == def {
			r.FileName = defaultFile
		}
		//Parse Search parameters
		r.SearchBody = strings.Fields(after(request, where))

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
		r.FileName = strings.TrimSpace(after(request, from))
		if r.FileName == "" {
			err := fmt.Errorf("no file name")
			return err
		}
		if r.FileName == def {
			r.FileName = defaultFile
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
				str := strings.TrimSpace(between(strings.Join(searchBody, " "), searchParam[i], and))
				searchValue = append(searchValue, str)
				for j := range searchBody {
					if j < len(searchBody)-1 && !newParam && searchBody[j] == and {
						searchBody = searchBody[j+1:]
						newParam = true
					}
				}
			} else {
				str := strings.TrimSpace(after(strings.Join(searchBody, " "), searchParam[i]))
				searchValue = append(searchValue, str)
			}
		}
	} else {
		str := strings.TrimSpace(after(strings.Join(searchBody, " "), searchParam[0]))
		searchValue = append(searchValue, str)
	}
	if searchValue == nil {
		err := fmt.Errorf("no search values")
		return nil, err
	}
	return searchValue, nil
}
