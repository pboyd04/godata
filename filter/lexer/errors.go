package lexer

import (
	"fmt"
)

type NoMatchingTokenError struct {
	Position int
}

type UnsupportedReplacementError struct {
	message string
}

func (e NoMatchingTokenError) Error() string {
	return fmt.Sprintf("no matching token at position %d", e.Position)
}

func newNoMatchingTokenError(position int) error {
	return NoMatchingTokenError{Position: position}
}

func (e UnsupportedReplacementError) Error() string {
	return e.message
}

func newUnsupportedReplacementError(format string, a ...interface{}) error {
	return &UnsupportedReplacementError{message: fmt.Sprintf(format, a...)}
}
