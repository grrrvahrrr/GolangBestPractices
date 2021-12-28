package process

import (
	"strconv"
	"strings"
)

func processSearchParam(searchParam string, indexParam []int, rec []string, v int, i int, searchValue []string) (string, error) {
	const equal string = "="
	const notEqual string = "!="
	const more string = ">"
	const moreEqual string = ">="
	const less string = "<"
	const lessEqual string = "<="
	const per string = `"`

	if strings.Contains(searchValue[i], per) || searchValue[i] == "true" || searchValue[i] == "false" {
		switch searchParam {
		case equal:
			if strings.Trim(rec[v], per) == strings.Trim(searchValue[i], per) {
				return rec[indexParam[i]], nil
			}
		case notEqual:
			if strings.Trim(rec[v], per) != strings.Trim(searchValue[i], per) {
				return rec[indexParam[i]], nil
			}
		}

	} else {
		recValue, err := strconv.ParseFloat(rec[v], 64)
		if err != nil {
			return "", err
		}
		value, err := strconv.ParseFloat(searchValue[i], 64)
		if err != nil {
			return "", err
		}

		switch searchParam {
		case more:
			if recValue > value {
				return rec[indexParam[i]], nil
			}
		case moreEqual:
			if recValue >= value {
				return rec[indexParam[i]], nil
			}
		case less:
			if recValue < value {
				return rec[indexParam[i]], nil
			}
		case lessEqual:
			if recValue <= value {
				return rec[indexParam[i]], nil
			}
		case equal:
			if recValue == value {
				return rec[indexParam[i]], nil
			}
		case notEqual:
			if recValue != value {
				return rec[indexParam[i]], nil
			}
		}
	}
	return "", nil
}
