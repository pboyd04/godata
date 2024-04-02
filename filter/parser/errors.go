package parser

import "fmt"

type ParsingError struct {
	message string
}

func (e *ParsingError) Error() string {
	return e.message
}

func newParserError(format string, a ...interface{}) error {
	return &ParsingError{message: fmt.Sprintf(format, a...)}
}

var ErrMoreThanOneChild = newParserError("more than one child")
var ErrNoSuchLanguage = newParserError("no such language")
