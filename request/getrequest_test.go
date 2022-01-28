package request

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	defaultFile := "/testDir/config.env"
	testRequest := "SELECT some test user request FROM some file\n"

	input := []byte(testRequest)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	//No errors found
	testString, err := GetRequest(defaultFile)
	if !assert.Equal(t, testRequest, testString) {
		t.Error("Wrong output by GetRequest func")
	}
	if !assert.Nil(t, err) {
		t.Error("Error in GetRequest func")
	}

}

func TestWrongRequest(t *testing.T) {
	defaultFile := "/testDir/config.env"
	testRequest := "some wrong user request\n"

	input := []byte(testRequest)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	//No errors found
	testString, err := GetRequest(defaultFile)
	if !assert.NotEqual(t, testRequest, testString) {
		t.Error("Wrong output by GetRequest func")
	}
	if !assert.NotNil(t, err) {
		t.Error("Should be an error in GetRequest func")
	}

}
