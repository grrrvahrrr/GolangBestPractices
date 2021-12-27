package process

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func processSearchParam(searchParam string, indexParam []int, rec []string, v int, i int, searchValue []string) string {
	const equal string = "="
	const notEqual string = "!="
	if strings.Contains(searchValue[i], `"`) || searchValue[i] == "true" || searchValue[i] == "false" {
		switch searchParam {
		case equal:
			if strings.Trim(rec[v], `"`) == strings.Trim(searchValue[i], `"`) {
				return rec[indexParam[i]]
			}
		case notEqual:
			if strings.Trim(rec[v], `"`) != strings.Trim(searchValue[i], `"`) {
				return rec[indexParam[i]]
			}
		}

	} else {
		recValue, err := strconv.ParseFloat(rec[v], 64)
		if err != nil {
			log.WithError(err).Debug("value wasn't a float")
		}
		value, err := strconv.ParseFloat(searchValue[i], 64)
		if err != nil {
			log.WithError(err).Debug("value wasn't a float")
		}

		switch searchParam {
		case ">":
			if recValue > value {
				return rec[indexParam[i]]
			}
		case ">=":
			if recValue >= value {
				return rec[indexParam[i]]
			}
		case "<":
			if recValue < value {
				return rec[indexParam[i]]
			}
		case "<=":
			if recValue <= value {
				return rec[indexParam[i]]
			}
		case equal:
			if recValue == value {
				return rec[indexParam[i]]
			}
		case notEqual:
			if recValue != value {
				return rec[indexParam[i]]
			}
		}
	}
	return ""
}
