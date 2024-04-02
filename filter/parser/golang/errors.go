package golang

import (
	"fmt"

	"github.com/pboyd04/godata/filter/lexer"
)

type UnsupportedOperandError struct {
	operand interface{}
}

type UnsupportedOperatorError struct {
	operator lexer.TokenKey
}

type UnsupportedDataTypeError struct {
}

type UnknownFieldError struct {
	field string
}

type ParserError struct {
	message string
}

func (e *UnsupportedOperandError) Error() string {
	return fmt.Sprintf("unsupported operand: %#v", e.operand)
}

func newUnsupportedOperandError(operand interface{}) error {
	return &UnsupportedOperandError{operand: operand}
}

func (e *UnsupportedDataTypeError) Error() string {
	return "unsupported data type"
}

func (e *ParserError) Error() string {
	return e.message
}

func newParserError(message string) error {
	return &ParserError{message: message}
}

func (e *UnknownFieldError) Error() string {
	return "unknown field: " + e.field
}

func (e *UnsupportedOperatorError) Error() string {
	return "unsupported operator: " + e.operator.String()
}
