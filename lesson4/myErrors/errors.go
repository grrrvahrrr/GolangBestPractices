package myErrors

import (
	"github.com/pkg/errors"
)

type myError struct {
	text string
}

func (e *myError) Error() string {
	return e.text
}

func CreateError(text string) error {
	return &myError{text: text}
}

func CheckError(text string) error {
	err := CreateError(text)
	if err != nil {
		return errors.Wrap(err, "Error")
	}
	return nil
}
