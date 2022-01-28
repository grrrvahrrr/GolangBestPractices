package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedRequest struct {
	FileName        string
	ColumnName      []string
	SearchBody      []string
	SearchParamName []string
	SearchParam     []string
	SearchValue     []string
}

func TestParseRequest(t *testing.T) {
	request := `SELECT SNo, Province/State FROM default WHERE Country/Region = "Mainland China" AND Confirmed > 100 AND Deaths < 50 AND Recovered > 20`
	defaultFile := "testdir/test.csv"

	var r Request
	var mr MockedRequest

	mr.FileName = "testdir/test.csv"
	mr.ColumnName = []string{"SNo", "Province/State"}
	mr.SearchBody = []string{"Country/Region", "=", `"Mainland`, `China"`, "AND", "Confirmed", ">", "100", "AND", "Deaths", "<", "50", "AND", "Recovered", ">", "20"}
	mr.SearchParamName = []string{"Country/Region", "Confirmed", "Deaths", "Recovered"}
	mr.SearchParam = []string{"=", ">", "<", ">"}
	mr.SearchValue = []string{`"Mainland China"`, "100", "50", "20"}

	err := r.ParseRequest(request, defaultFile)
	if !assert.Nil(t, err) {
		t.Error("ParseRequest returned an error")
	}

	if !assert.Equal(t, mr.FileName, r.FileName) {
		t.Error("Filename is not equal to expected")
	}

	for i, v := range r.ColumnName {
		if !assert.Equal(t, mr.ColumnName[i], v) {
			t.Error("ColumnName is not equal to expected")
		}
	}

	for i, v := range r.SearchBody {
		if !assert.Equal(t, mr.SearchBody[i], v) {
			t.Error("SearchBody is not equal to expected")
		}
	}

	for i, v := range r.SearchParamName {
		if !assert.Equal(t, mr.SearchParamName[i], v) {
			t.Error("SearchParamName is not equal to expected")
		}
	}

	for i, v := range r.SearchParam {
		if !assert.Equal(t, mr.SearchParam[i], v) {
			t.Error("SearchParam is not equal to expected")
		}
	}

	for i, v := range r.SearchValue {
		if !assert.Equal(t, mr.SearchValue[i], v) {
			t.Error("SearchValue is not equal to expected")
		}
	}

}

func TestBetween(t *testing.T) {
	testString := "I test between function."
	expectedString := " between "

	actualString := between(testString, "test", "function")
	if !assert.Equal(t, expectedString, actualString) {
		t.Error("Between() failed")
	}
}

func TestAfter(t *testing.T) {
	testString := "I test between function."
	expectedString := " between function."

	actualString := after(testString, "test")
	if !assert.Equal(t, expectedString, actualString) {
		t.Error("After() failed")
	}
}

func TestParseSearchValue(t *testing.T) {
	searchBody := []string{"Country/Region", "=", `"Mainland China"`, "AND", "Confirmed", ">", "100", "AND", "Deaths", "<", "50", "AND", "Recovered", ">", "20"}
	searchParam := []string{"=", ">", "<", ">"}

	expectedSearchValue := []string{`"Mainland China"`, "100", "50", "20"}

	actualSearchValue, err := parseSearchValue(searchBody, searchParam)
	if !assert.Nil(t, err) {
		t.Error("parseSearchValue returned an error")
	}

	for i, v := range expectedSearchValue {
		if !assert.Equal(t, actualSearchValue[i], v) {
			t.Error("SearchValue is not equal to expected")
		}
	}

}
