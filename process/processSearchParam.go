package process

import (
	"strconv"
)

func processMore(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue > value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}

func processMoreEqual(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue >= value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}

func processLess(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue < value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}

func processLessEqual(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue <= value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}

func processEqual(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue == value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}

func processNotEqual(indexParam []int, rec []string, v int, i int, searchValue []string) (stringToAdd string) {
	recValue, _ := strconv.ParseFloat(rec[v], 64)
	value, _ := strconv.ParseFloat(searchValue[i], 64)
	if recValue != value {
		stringToAdd = rec[indexParam[i]]
		return stringToAdd
	}
	return ""
}
