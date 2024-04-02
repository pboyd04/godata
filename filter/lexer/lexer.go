package lexer

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type TokenKey int

const (
	TokenTrue = iota
	TokenFalse
	UnquotedString
	SingleQuotedString
	DoubleQuotedString
	OpenParens
	CloseParens
	OpenSquareBracket
	CloseSquareBracket
	OpenCurlyBrace
	CloseCurlyBrace
	Colon
	Not
	Equals
	NotEquals
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	Has
	In
	Concat
	Contains
	EndsWith
	IndexOf
	Length
	StartsWith
	Substring
	HasSubset
	HasSubsequence
	MatchesPattern
	ToLower
	ToUpper
	Trim
	Day
	FractionalSeconds
	Hour
	Minute
	Month
	Second
	Year
	Ceiling
	Floor
	Round
	Add
	Subtract
	Multiply
	Divide
	DivideFloat
	Modulo
	And
	Or
	NullLiteral
	FloatingPointLiteral
	IntegerLiteral
	Comma
	// Not currently supported: date, maxDateTime, minDateTime, now, time, totalOffsetMinutes, totalSeconds, cast, isOf, geo.*, case, any, all.
)

type tokenMatcher func(string, *Lexer) int

type TokenType struct {
	typeKey     TokenKey
	regex       *regexp.Regexp
	stringMatch *string
	matcherFn   tokenMatcher
}

type Lexer struct {
	text     string
	lower    string
	position int
	length   int
	types    []TokenType
}

type Token struct {
	Type  TokenKey
	Start int
	End   int
	Text  string
}

func ptrFromConst(s string) *string {
	return &s
}

func singleQuoteString(s string, _ *Lexer) int {
	if s[0] == '\'' {
		length := len(s)
		for i := 1; i < length; i++ {
			if s[i] == '\'' {
				return i + 1
			}
		}
	}
	return -1
}

func doubleQuoteString(s string, _ *Lexer) int {
	if s[0] == '\x22' { // double quote
		length := len(s)
		for i := 1; i < length; i++ {
			if s[i] == '\x22' {
				return i + 1
			}
		}
	}
	return -1
}

// Replaces regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+`) as this is faster.
//
//nolint:cyclop // This is a simple function that is easy to understand
func testForFloat(s string, _ *Lexer) int {
	startIndex := 0
	if s[0] == '-' || s[0] == '+' {
		startIndex = 1
	}
	length := len(s)
	if startIndex >= length || !unicode.IsDigit(rune(s[startIndex])) {
		return -1
	}
	foundDot := false
	for i := startIndex + 1; i < length; i++ {
		if s[i] == '.' {
			if foundDot {
				return i
			}
			foundDot = true
			continue
		}
		if !unicode.IsDigit(rune(s[i])) {
			if !foundDot {
				return -1
			}
			return i
		}
	}
	if !foundDot {
		return -1
	}
	return length
}

// Replaces regexp.MustCompile(`^[-+]?[0-9]+\b`) as this is faster.
//
//nolint:cyclop // This is a simple function that is easy to understand
func testForInt(s string, _ *Lexer) int {
	startIndex := 0
	if s[0] == '-' || s[0] == '+' {
		startIndex = 1
	}
	length := len(s)
	if startIndex >= length || !unicode.IsDigit(rune(s[startIndex])) {
		return -1
	}
	for i := startIndex + 1; i < length; i++ {
		if unicode.IsSpace(rune(s[i])) || s[i] == ',' || s[i] == ')' || s[i] == ']' || s[i] == '}' {
			return i
		}
		if !unicode.IsDigit(rune(s[i-1])) {
			// Found a letter right next to a digit this is something like a mongo id, treat it like a string
			return -1
		}
	}
	return length
}

// Replaces regexp.MustCompile(`([^\s,)'"]+)`) as this is faster.
func testForUnquotedString(s string, _ *Lexer) int {
	length := len(s)
	for i := 0; i < length; i++ {
		if unicode.IsSpace(rune(s[i])) || s[i] == ',' || s[i] == ')' || s[i] == ']' || s[i] == '}' || s[i] == '\'' || s[i] == '"' {
			return i
		}
	}
	return length
}

//nolint:gochecknoglobals // We only need to perform all this init once, otherwise we pay it every time we lex a string
var odataLexTypes = []TokenType{
	{TokenTrue, nil, ptrFromConst("true"), nil},
	{TokenFalse, nil, ptrFromConst("false"), nil},
	{SingleQuotedString, nil, nil, singleQuoteString},
	{DoubleQuotedString, nil, nil, doubleQuoteString},
	{OpenParens, nil, ptrFromConst("("), nil},
	{CloseParens, nil, ptrFromConst(")"), nil},
	{OpenSquareBracket, nil, ptrFromConst("["), nil},
	{CloseSquareBracket, nil, ptrFromConst("]"), nil},
	{OpenCurlyBrace, nil, ptrFromConst("{"), nil},
	{CloseCurlyBrace, nil, ptrFromConst("}"), nil},
	{Colon, nil, ptrFromConst(":"), nil},
	{Equals, nil, ptrFromConst("eq "), nil},
	{NotEquals, nil, ptrFromConst("ne "), nil},
	{GreaterThan, nil, ptrFromConst("gt "), nil},
	{GreaterThanOrEqual, nil, ptrFromConst("ge "), nil},
	{LessThan, nil, ptrFromConst("lt "), nil},
	{LessThanOrEqual, nil, ptrFromConst("le "), nil},
	{And, nil, ptrFromConst("and "), nil},
	{Or, nil, ptrFromConst("or "), nil},
	{Not, nil, ptrFromConst("not "), nil},
	{Has, nil, ptrFromConst("has "), nil},
	{In, nil, ptrFromConst("in "), nil},
	{Concat, nil, ptrFromConst("concat"), nil},
	{Contains, nil, ptrFromConst("contains"), nil},
	{EndsWith, nil, ptrFromConst("endswith"), nil},
	{IndexOf, nil, ptrFromConst("indexof"), nil},
	{Length, nil, ptrFromConst("length"), nil},
	{StartsWith, nil, ptrFromConst("startswith"), nil},
	{Substring, nil, ptrFromConst("substring"), nil},
	{HasSubset, nil, ptrFromConst("hassubset"), nil},
	{HasSubsequence, nil, ptrFromConst("hassubsequence"), nil},
	{MatchesPattern, nil, ptrFromConst("matchespattern"), nil},
	{ToLower, nil, ptrFromConst("tolower"), nil},
	{ToUpper, nil, ptrFromConst("toupper"), nil},
	{Trim, nil, ptrFromConst("trim"), nil},
	{Day, nil, ptrFromConst("day"), nil},
	{FractionalSeconds, nil, ptrFromConst("fractionalseconds"), nil},
	{Hour, nil, ptrFromConst("hour"), nil},
	{Minute, nil, ptrFromConst("minute"), nil},
	{Month, nil, ptrFromConst("month"), nil},
	{Second, nil, ptrFromConst("second"), nil},
	{Year, nil, ptrFromConst("year"), nil},
	{Ceiling, nil, ptrFromConst("ceiling"), nil},
	{Floor, nil, ptrFromConst("floor"), nil},
	{Round, nil, ptrFromConst("round"), nil},
	{Add, nil, ptrFromConst("add "), nil},
	{Subtract, nil, ptrFromConst("sub "), nil},
	{Multiply, nil, ptrFromConst("mul "), nil},
	{Divide, nil, ptrFromConst("div "), nil},
	{DivideFloat, nil, ptrFromConst("divby "), nil},
	{Modulo, nil, ptrFromConst("mod "), nil},
	{NullLiteral, nil, ptrFromConst("null"), nil},
	{Comma, nil, ptrFromConst(","), nil},
	{FloatingPointLiteral, nil, nil, testForFloat},
	{IntegerLiteral, nil, nil, testForInt},
	// Needs to be near the end otherwise it will match everything
	{UnquotedString, nil, nil, testForUnquotedString},
}

//nolint:funlen
func NewLexer(input string) *Lexer {
	ret := &Lexer{text: input, position: 0}
	// Avoid the need for case-insensitive regex/string compare
	ret.lower = strings.ToLower(input)
	ret.types = odataLexTypes
	ret.length = len(input)
	if ret.length != len(ret.lower) {
		ret.length = -1
	}
	return ret
}

func (l *Lexer) testMatchingFunction(t TokenType) (*Token, error) {
	length := t.matcherFn(l.lower[l.position:], l)
	if length > 0 {
		start := l.position
		if start+length > l.length {
			return nil, newNoMatchingTokenError(start)
		}
		l.position += length
		return &Token{Type: t.typeKey, Start: start, End: l.position, Text: l.text[start:l.position]}, nil
	}
	return nil, nil
}

func (l *Lexer) testStringMatch(t TokenType) (*Token, error) {
	if strings.HasPrefix(l.lower[l.position:], *t.stringMatch) {
		if t.typeKey.HasParameters() {
			// If the next character is a ( then we need to return the function name and the open parens
			length := len(*t.stringMatch)
			if l.position+length < len(l.text) && l.text[l.position+length] == '(' {
				start := l.position
				l.position += length // Don't return the ( as part of the token
				return &Token{Type: t.typeKey, Start: start, End: l.position, Text: l.text[start:l.position]}, nil
			}
			return nil, nil
		}
		start := l.position
		length := len(*t.stringMatch)
		if start+length > l.length {
			return nil, newNoMatchingTokenError(start)
		}
		l.position += length
		return &Token{Type: t.typeKey, Start: start, End: l.position, Text: l.text[start:l.position]}, nil
	}
	return nil, nil
}

func (l *Lexer) testRegex(t TokenType) *Token {
	ret := t.regex.Find([]byte(l.lower[l.position:]))
	if ret != nil {
		start := l.position
		l.position += len(ret)
		return &Token{Type: t.typeKey, Start: start, End: l.position, Text: l.text[start:l.position]}
	}
	return nil
}

func (l *Lexer) testToken(t TokenType) (*Token, error) {
	if t.matcherFn != nil {
		res, err := l.testMatchingFunction(t)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	if t.stringMatch != nil && strings.HasPrefix(l.lower[l.position:], *t.stringMatch) {
		res, err := l.testStringMatch(t)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	if t.regex != nil {
		res := l.testRegex(t)
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}

func (l *Lexer) NextToken() (*Token, error) {
	if l.position >= l.length {
		return nil, nil
	}
	// Skip whitespace
	if unicode.IsSpace(rune(l.text[l.position])) {
		l.position++
		if l.position >= l.length {
			return nil, nil
		}
	}
	for _, t := range l.types {
		res, err := l.testToken(t)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, newNoMatchingTokenError(l.position)
}

func (t *Token) IsUnary() bool {
	return t.Type == Not
}

func (t *Token) IsMultiplicative() bool {
	return t.Type == Multiply || t.Type == Divide || t.Type == DivideFloat || t.Type == Modulo
}

func (t *Token) IsAdditive() bool {
	return t.Type == Add || t.Type == Subtract
}

func (t *Token) IsRelational() bool {
	return t.Type == GreaterThan || t.Type == GreaterThanOrEqual || t.Type == LessThan || t.Type == LessThanOrEqual
}

func (t *Token) IsEquality() bool {
	return t.Type == Equals || t.Type == NotEquals
}

func (t TokenKey) HasParameters() bool {
	switch t {
	case Concat, Contains, EndsWith, IndexOf, Length, StartsWith, Substring, HasSubset, HasSubsequence,
		MatchesPattern, ToLower, ToUpper, Trim, Day, FractionalSeconds, Hour, Minute, Month, Second,
		Year, Ceiling, Floor, Round:
		return true
	default:
		return false
	}
}

func (t *Token) HasParameters() bool {
	return t.Type.HasParameters()
}

func (t Token) GetData() (interface{}, error) {
	str := t.Text
	switch t.Type {
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	case SingleQuotedString, DoubleQuotedString:
		// Remove the quotes
		return str[1 : len(str)-1], nil
	case NullLiteral:
		return nil, nil
	case FloatingPointLiteral:
		return strconv.ParseFloat(str, 64)
	case IntegerLiteral:
		return strconv.Atoi(str)
	default:
		return str, nil
	}
}

func (t *Token) IsCorrectReplacement(index int) bool {
	var strCmp string
	switch t.Type {
	case SingleQuotedString:
		strCmp = "':" + strconv.Itoa(index) + "'"
	case DoubleQuotedString:
		strCmp = `":` + strconv.Itoa(index) + `"`
	default:
		return false
	}
	return t.Text == strCmp
}

func (t *Token) Replace(operand interface{}) error {
	switch operand := operand.(type) {
	case string:
		t.Text = "'" + operand + "'"
	case int:
		t.Text = strconv.Itoa(operand)
		t.Type = IntegerLiteral
	case float64:
		t.Text = strconv.FormatFloat(operand, 'f', -1, 64)
		t.Type = FloatingPointLiteral
	default:
		return newUnsupportedReplacementError("unsupported type %T", operand)
	}
	return nil
}

func (t *TokenKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

//nolint:funlen,gocyclo,cyclop
func (t TokenKey) String() string {
	switch t {
	case TokenTrue:
		return "TokenTrue"
	case TokenFalse:
		return "TokenFalse"
	case UnquotedString:
		return "UnquotedString"
	case SingleQuotedString:
		return "SingleQuotedString"
	case DoubleQuotedString:
		return "DoubleQuotedString"
	case OpenParens:
		return "OpenParens"
	case CloseParens:
		return "CloseParens"
	case OpenSquareBracket:
		return "OpenSquareBracket"
	case CloseSquareBracket:
		return "CloseSquareBracket"
	case OpenCurlyBrace:
		return "OpenCurlyBrace"
	case CloseCurlyBrace:
		return "CloseCurlyBrace"
	case Colon:
		return "Colon"
	case Not:
		return "Not"
	case Equals:
		return "Equals"
	case NotEquals:
		return "NotEquals"
	case GreaterThan:
		return "GreaterThan"
	case GreaterThanOrEqual:
		return "GreaterThanOrEqual"
	case LessThan:
		return "LessThan"
	case LessThanOrEqual:
		return "LessThanOrEqual"
	case Has:
		return "Has"
	case In:
		return "In"
	case Concat:
		return "Concat"
	case Contains:
		return "Contains"
	case EndsWith:
		return "EndsWith"
	case IndexOf:
		return "IndexOf"
	case Length:
		return "Length"
	case StartsWith:
		return "StartsWith"
	case Substring:
		return "Substring"
	case HasSubset:
		return "HasSubset"
	case HasSubsequence:
		return "HasSubsequence"
	case MatchesPattern:
		return "MatchesPattern"
	case ToLower:
		return "ToLower"
	case ToUpper:
		return "ToUpper"
	case Trim:
		return "Trim"
	case Day:
		return "Day"
	case FractionalSeconds:
		return "FractionalSeconds"
	case Hour:
		return "Hour"
	case Minute:
		return "Minute"
	case Month:
		return "Month"
	case Second:
		return "Second"
	case Year:
		return "Year"
	case Ceiling:
		return "Ceiling"
	case Floor:
		return "Floor"
	case Round:
		return "Round"
	case Add:
		return "Add"
	case Subtract:
		return "Subtract"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case DivideFloat:
		return "DivideFloat"
	case Modulo:
		return "Modulo"
	case And:
		return "And"
	case Or:
		return "Or"
	case NullLiteral:
		return "NullLiteral"
	case FloatingPointLiteral:
		return "FloatingPointLiteral"
	case IntegerLiteral:
		return "IntegerLiteral"
	case Comma:
		return "Comma"
	default:
		return strconv.Itoa(int(t))
	}
}
