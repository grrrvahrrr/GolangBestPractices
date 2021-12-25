package process

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

func processMore(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	value, err := strconv.ParseFloat(searchValue[i], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	if recValue > value {
		return rec[indexParam[i]]
	}
	return ""
}

func processMoreEqual(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	value, err := strconv.ParseFloat(searchValue[i], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	if recValue >= value {
		return rec[indexParam[i]]
	}
	return ""
}

func processLess(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	value, err := strconv.ParseFloat(searchValue[i], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	if recValue < value {
		return rec[indexParam[i]]
	}
	return ""
}

func processLessEqual(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	value, err := strconv.ParseFloat(searchValue[i], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
	}
	if recValue <= value {
		return rec[indexParam[i]]
	}
	return ""
}

func processEqual(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
		if rec[v] == searchValue[i] {
			return rec[indexParam[i]]
		}
	} else {
		value, err := strconv.ParseFloat(searchValue[i], 64)
		if err != nil {
			log.WithError(err).Debug("value wasn't a float")
		}
		if recValue == value {
			return rec[indexParam[i]]
		}
	}
	return ""
}

func processNotEqual(indexParam []int, rec []string, v int, i int, searchValue []string) string {
	recValue, err := strconv.ParseFloat(rec[v], 64)
	if err != nil {
		log.WithError(err).Debug("value wasn't a float")
		if rec[v] != searchValue[i] {
			return rec[indexParam[i]]
		}
	} else {
		value, err := strconv.ParseFloat(searchValue[i], 64)
		if err != nil {
			log.WithError(err).Debug("value wasn't a float")
		}
		if recValue != value {
			return rec[indexParam[i]]
		}
	}

	return ""
}
