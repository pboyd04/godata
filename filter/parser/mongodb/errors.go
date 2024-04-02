package mongodb

import (
	"fmt"

	"github.com/pboyd04/godata/filter/parser"
)

type UnsupportedOperandError struct {
	operand interface{}
}

type UnsupportedOperatorError struct {
	operator parser.Operator
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

func (e *UnsupportedOperatorError) Error() string {
	return fmt.Sprintf("unsupported operator: %v", e.operator)
}

func newUnsupportedOperatorError(operator parser.Operator) error {
	return &UnsupportedOperatorError{operator: operator}
}

func (e *ParserError) Error() string {
	return e.message
}

func newParserError(message string) error {
	return &ParserError{message: message}
}
