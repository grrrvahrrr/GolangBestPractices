package main

import (
	"github.com/pkg/errors"
)

type myError struct {
	text string
}

func (e *myError) Error() string {
	return e.text
}

func createError(text string) error {
	return &myError{text: text}
}

func checkError(text string) error {
	err := createError(text)
	if err != nil {
		return errors.Wrap(err, "Error")
	}
	return nil
}
