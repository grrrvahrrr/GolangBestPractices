package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessSearchParam(t *testing.T) {
	searchParam := []string{"="}
	indexParam := []int{3}
	rec := []string{"306427", "date", "Zhejiang", `"Mainland China"`, "1364.0", "1.0", "1324.0"}
	i := 0
	searchValue := []string{`"Mainland China"`}

	expectedString := `"Mainland China"`

	actualString, err := processSearchParam(searchParam[i], indexParam, rec, i, searchValue)
	if !assert.Nil(t, err) {
		t.Error("processSearchParam returned an error")
	}
	if !assert.Equal(t, expectedString, actualString) {
		t.Error("output is not equal to expected")
	}
}
