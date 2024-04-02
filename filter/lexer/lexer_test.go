package lexer_test

import (
	"testing"

	"github.com/pboyd04/godata/filter/lexer"
)

type testData struct {
	input    string
	expected []lexer.Token
}

//nolint:gochecknoglobals,dupl // Just test case data
var tests = []testData{
	{
		input: "true",
		expected: []lexer.Token{
			{Type: lexer.TokenTrue, Start: 0, End: 4},
		},
	},
	{
		input: "false",
		expected: []lexer.Token{
			{Type: lexer.TokenFalse, Start: 0, End: 5},
		},
	},
	{
		input: "Name eq 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "(Name eq 'Milk')",
		expected: []lexer.Token{
			{Type: lexer.OpenParens, Start: 0, End: 1},
			{Type: lexer.UnquotedString, Start: 1, End: 5},
			{Type: lexer.Equals, Start: 6, End: 9},
			{Type: lexer.SingleQuotedString, Start: 9, End: 15},
			{Type: lexer.CloseParens, Start: 15, End: 16},
		},
	},
	{
		input: "Name ne 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.NotEquals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "Name gt 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.GreaterThan, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "Name ge 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.GreaterThanOrEqual, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "Name lt 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.LessThan, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "Name le 'Milk'",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.LessThanOrEqual, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
		},
	},
	{
		input: "Name eq 'Milk' and Price lt 2.55",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
			{Type: lexer.And, Start: 15, End: 19},
			{Type: lexer.UnquotedString, Start: 19, End: 24},
			{Type: lexer.LessThan, Start: 25, End: 28},
			{Type: lexer.FloatingPointLiteral, Start: 28, End: 32},
		},
	},
	{
		input: "Name EQ 'Milk' AND Price LT 2.55",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
			{Type: lexer.And, Start: 15, End: 19},
			{Type: lexer.UnquotedString, Start: 19, End: 24},
			{Type: lexer.LessThan, Start: 25, End: 28},
			{Type: lexer.FloatingPointLiteral, Start: 28, End: 32},
		},
	},
	{
		input: "Name eq 'Milk' AND Price LT 2.55",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
			{Type: lexer.And, Start: 15, End: 19},
			{Type: lexer.UnquotedString, Start: 19, End: 24},
			{Type: lexer.LessThan, Start: 25, End: 28},
			{Type: lexer.FloatingPointLiteral, Start: 28, End: 32},
		},
	},
	{
		input: "Name eq 'Milk' or Price lt 2.55",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 14},
			{Type: lexer.Or, Start: 15, End: 18},
			{Type: lexer.UnquotedString, Start: 18, End: 23},
			{Type: lexer.LessThan, Start: 24, End: 27},
			{Type: lexer.FloatingPointLiteral, Start: 27, End: 31},
		},
	},
	{
		input: "Name in ('Milk', 'Cheese')",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.In, Start: 5, End: 8},
			{Type: lexer.OpenParens, Start: 8, End: 9},
			{Type: lexer.SingleQuotedString, Start: 9, End: 15},
			{Type: lexer.Comma, Start: 15, End: 16},
			{Type: lexer.SingleQuotedString, Start: 17, End: 25},
			{Type: lexer.CloseParens, Start: 25, End: 26},
		},
	},
	{
		input: "Name in ['Milk', 'Cheese']",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.In, Start: 5, End: 8},
			{Type: lexer.OpenSquareBracket, Start: 8, End: 9},
			{Type: lexer.SingleQuotedString, Start: 9, End: 15},
			{Type: lexer.Comma, Start: 15, End: 16},
			{Type: lexer.SingleQuotedString, Start: 17, End: 25},
			{Type: lexer.CloseSquareBracket, Start: 25, End: 26},
		},
	},
	{
		input: "_id eq 6206b158000e1859781d5e16",
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 3},
			{Type: lexer.Equals, Start: 4, End: 7},
			{Type: lexer.UnquotedString, Start: 7, End: 31},
		},
	},
	{
		input: "contains(Name,'red')",
		expected: []lexer.Token{
			{Type: lexer.Contains, Start: 0, End: 8},
			{Type: lexer.OpenParens, Start: 8, End: 9},
			{Type: lexer.UnquotedString, Start: 9, End: 13},
			{Type: lexer.Comma, Start: 13, End: 14},
			{Type: lexer.SingleQuotedString, Start: 14, End: 19},
			{Type: lexer.CloseParens, Start: 19, End: 20},
		},
	},
	{

		input: `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 7},
			{Type: lexer.Equals, Start: 8, End: 11},
			{Type: lexer.OpenCurlyBrace, Start: 11, End: 12},
			{Type: lexer.DoubleQuotedString, Start: 12, End: 20},
			{Type: lexer.Colon, Start: 20, End: 21},
			{Type: lexer.DoubleQuotedString, Start: 21, End: 30},
			{Type: lexer.Comma, Start: 30, End: 31},
			{Type: lexer.DoubleQuotedString, Start: 31, End: 37},
			{Type: lexer.Colon, Start: 37, End: 38},
			{Type: lexer.DoubleQuotedString, Start: 38, End: 47},
			{Type: lexer.Comma, Start: 47, End: 48},
			{Type: lexer.DoubleQuotedString, Start: 48, End: 55},
			{Type: lexer.Colon, Start: 55, End: 56},
			{Type: lexer.DoubleQuotedString, Start: 56, End: 60},
			{Type: lexer.Comma, Start: 60, End: 61},
			{Type: lexer.DoubleQuotedString, Start: 61, End: 70},
			{Type: lexer.Colon, Start: 70, End: 71},
			{Type: lexer.DoubleQuotedString, Start: 71, End: 78},
			{Type: lexer.CloseCurlyBrace, Start: 78, End: 79},
		},
	},
	{
		input: "endswith(Name,'ilk')",
		expected: []lexer.Token{
			{Type: lexer.EndsWith, Start: 0, End: 8},
			{Type: lexer.OpenParens, Start: 8, End: 9},
			{Type: lexer.UnquotedString, Start: 9, End: 13},
			{Type: lexer.Comma, Start: 13, End: 14},
			{Type: lexer.SingleQuotedString, Start: 14, End: 19},
			{Type: lexer.CloseParens, Start: 19, End: 20},
		},
	},
	{
		input: "not endswith(Name,'ilk')",
		expected: []lexer.Token{
			{Type: lexer.Not, Start: 0, End: 4},
			{Type: lexer.EndsWith, Start: 4, End: 12},
			{Type: lexer.OpenParens, Start: 12, End: 13},
			{Type: lexer.UnquotedString, Start: 13, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.SingleQuotedString, Start: 18, End: 23},
			{Type: lexer.CloseParens, Start: 23, End: 24},
		},
	},
	{
		input: "length(CompanyName) eq 19",
		expected: []lexer.Token{
			{Type: lexer.Length, Start: 0, End: 6},
			{Type: lexer.OpenParens, Start: 6, End: 7},
			{Type: lexer.UnquotedString, Start: 7, End: 18},
			{Type: lexer.CloseParens, Start: 18, End: 19},
			{Type: lexer.Equals, Start: 20, End: 23},
			{Type: lexer.IntegerLiteral, Start: 23, End: 25},
		},
	},
	{
		input: "startswith(CompanyName,'Futterkiste')",
		expected: []lexer.Token{
			{Type: lexer.StartsWith, Start: 0, End: 10},
			{Type: lexer.OpenParens, Start: 10, End: 11},
			{Type: lexer.UnquotedString, Start: 11, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.SingleQuotedString, Start: 23, End: 36},
			{Type: lexer.CloseParens, Start: 36, End: 37},
		},
	},
	{
		input: `hassubset(Names,["Milk", "Cheese"])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubset, Start: 0, End: 9},
			{Type: lexer.OpenParens, Start: 9, End: 10},
			{Type: lexer.UnquotedString, Start: 10, End: 15},
			{Type: lexer.Comma, Start: 15, End: 16},
			{Type: lexer.OpenSquareBracket, Start: 16, End: 17},
			{Type: lexer.DoubleQuotedString, Start: 17, End: 23},
			{Type: lexer.Comma, Start: 23, End: 24},
			{Type: lexer.DoubleQuotedString, Start: 25, End: 33},
			{Type: lexer.CloseSquareBracket, Start: 33, End: 34},
			{Type: lexer.CloseParens, Start: 34, End: 35},
		},
	},
	{
		input: `Price add 2.45 eq 5.00`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Add, Start: 6, End: 10},
			{Type: lexer.FloatingPointLiteral, Start: 10, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.FloatingPointLiteral, Start: 18, End: 22},
		},
	},
	{
		input: `Price sub 0.55 eq 2.00`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Subtract, Start: 6, End: 10},
			{Type: lexer.FloatingPointLiteral, Start: 10, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.FloatingPointLiteral, Start: 18, End: 22},
		},
	},
	{
		input: `Price mul 2.0 eq 5.10`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Multiply, Start: 6, End: 10},
			{Type: lexer.FloatingPointLiteral, Start: 10, End: 13},
			{Type: lexer.Equals, Start: 14, End: 17},
			{Type: lexer.FloatingPointLiteral, Start: 17, End: 21},
		},
	},
	{
		input: `Price div 2.55 eq 1`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Divide, Start: 6, End: 10},
			{Type: lexer.FloatingPointLiteral, Start: 10, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
		},
	},
	{
		input: `Price div 2 eq 2`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Divide, Start: 6, End: 10},
			{Type: lexer.IntegerLiteral, Start: 10, End: 11},
			{Type: lexer.Equals, Start: 12, End: 15},
			{Type: lexer.IntegerLiteral, Start: 15, End: 16},
		},
	},
	{
		input: `Price divby 2 eq 2.5`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.DivideFloat, Start: 6, End: 12},
			{Type: lexer.IntegerLiteral, Start: 12, End: 13},
			{Type: lexer.Equals, Start: 14, End: 17},
			{Type: lexer.FloatingPointLiteral, Start: 17, End: 20},
		},
	},
	{
		input: `Rating mod 5 eq 0`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 6},
			{Type: lexer.Modulo, Start: 7, End: 11},
			{Type: lexer.IntegerLiteral, Start: 11, End: 12},
			{Type: lexer.Equals, Start: 13, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
		},
	},
	{
		input: `style has Sales.Pattern'Yellow'`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Has, Start: 6, End: 10},
			{Type: lexer.UnquotedString, Start: 10, End: 23},
			{Type: lexer.SingleQuotedString, Start: 23, End: 31},
		},
	},
	{
		input: `(4 add 5) mod (4 sub 1) eq 0`,
		expected: []lexer.Token{
			{Type: lexer.OpenParens, Start: 0, End: 1},
			{Type: lexer.IntegerLiteral, Start: 1, End: 2},
			{Type: lexer.Add, Start: 3, End: 7},
			{Type: lexer.IntegerLiteral, Start: 7, End: 8},
			{Type: lexer.CloseParens, Start: 8, End: 9},
			{Type: lexer.Modulo, Start: 10, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.IntegerLiteral, Start: 15, End: 16},
			{Type: lexer.Subtract, Start: 17, End: 21},
			{Type: lexer.IntegerLiteral, Start: 21, End: 22},
			{Type: lexer.CloseParens, Start: 22, End: 23},
			{Type: lexer.Equals, Start: 24, End: 27},
			{Type: lexer.IntegerLiteral, Start: 27, End: 28},
		},
	},
	{
		input: `concat(concat(City,', '),Country) eq 'Berlin, Germany'`,
		expected: []lexer.Token{
			{Type: lexer.Concat, Start: 0, End: 6},
			{Type: lexer.OpenParens, Start: 6, End: 7},
			{Type: lexer.Concat, Start: 7, End: 13},
			{Type: lexer.OpenParens, Start: 13, End: 14},
			{Type: lexer.UnquotedString, Start: 14, End: 18},
			{Type: lexer.Comma, Start: 18, End: 19},
			{Type: lexer.SingleQuotedString, Start: 19, End: 23},
			{Type: lexer.CloseParens, Start: 23, End: 24},
			{Type: lexer.Comma, Start: 24, End: 25},
			{Type: lexer.UnquotedString, Start: 25, End: 32},
			{Type: lexer.CloseParens, Start: 32, End: 33},
			{Type: lexer.Equals, Start: 34, End: 37},
			{Type: lexer.SingleQuotedString, Start: 37, End: 54},
		},
	},
	{
		input: `indexof(CompanyName,'lfreds') eq 1`,
		expected: []lexer.Token{
			{Type: lexer.IndexOf, Start: 0, End: 7},
			{Type: lexer.OpenParens, Start: 7, End: 8},
			{Type: lexer.UnquotedString, Start: 8, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.SingleQuotedString, Start: 20, End: 28},
			{Type: lexer.CloseParens, Start: 28, End: 29},
			{Type: lexer.Equals, Start: 30, End: 33},
			{Type: lexer.IntegerLiteral, Start: 33, End: 34},
		},
	},
	{
		input: `substring(CompanyName,1) eq 'lfreds Futterkiste'`,
		expected: []lexer.Token{
			{Type: lexer.Substring, Start: 0, End: 9},
			{Type: lexer.OpenParens, Start: 9, End: 10},
			{Type: lexer.UnquotedString, Start: 10, End: 21},
			{Type: lexer.Comma, Start: 21, End: 22},
			{Type: lexer.IntegerLiteral, Start: 22, End: 23},
			{Type: lexer.CloseParens, Start: 23, End: 24},
			{Type: lexer.Equals, Start: 25, End: 28},
			{Type: lexer.SingleQuotedString, Start: 28, End: 48},
		},
	},
	{
		input: `substring(CompanyName,1,2) eq 'lf'`,
		expected: []lexer.Token{
			{Type: lexer.Substring, Start: 0, End: 9},
			{Type: lexer.OpenParens, Start: 9, End: 10},
			{Type: lexer.UnquotedString, Start: 10, End: 21},
			{Type: lexer.Comma, Start: 21, End: 22},
			{Type: lexer.IntegerLiteral, Start: 22, End: 23},
			{Type: lexer.Comma, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.CloseParens, Start: 25, End: 26},
			{Type: lexer.Equals, Start: 27, End: 30},
			{Type: lexer.SingleQuotedString, Start: 30, End: 34},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,1,3])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.CloseSquareBracket, Start: 21, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.OpenSquareBracket, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.Comma, Start: 27, End: 28},
			{Type: lexer.IntegerLiteral, Start: 28, End: 29},
			{Type: lexer.CloseSquareBracket, Start: 29, End: 30},
			{Type: lexer.CloseParens, Start: 30, End: 31},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,1])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.CloseSquareBracket, Start: 21, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.OpenSquareBracket, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.CloseSquareBracket, Start: 27, End: 28},
			{Type: lexer.CloseParens, Start: 28, End: 29},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,3])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.CloseSquareBracket, Start: 21, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.OpenSquareBracket, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.CloseSquareBracket, Start: 27, End: 28},
			{Type: lexer.CloseParens, Start: 28, End: 29},
		},
	},
	{
		input: `hassubsequence([4,1,3,1],[1,1])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.Comma, Start: 21, End: 22},
			{Type: lexer.IntegerLiteral, Start: 22, End: 23},
			{Type: lexer.CloseSquareBracket, Start: 23, End: 24},
			{Type: lexer.Comma, Start: 24, End: 25},
			{Type: lexer.OpenSquareBracket, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.Comma, Start: 27, End: 28},
			{Type: lexer.IntegerLiteral, Start: 28, End: 29},
			{Type: lexer.CloseSquareBracket, Start: 29, End: 30},
			{Type: lexer.CloseParens, Start: 30, End: 31},
		},
	},
	{
		input: `hassubsequence([4,1,3],[1,3,4])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.CloseSquareBracket, Start: 21, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.OpenSquareBracket, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.Comma, Start: 27, End: 28},
			{Type: lexer.IntegerLiteral, Start: 28, End: 29},
			{Type: lexer.CloseSquareBracket, Start: 29, End: 30},
			{Type: lexer.CloseParens, Start: 30, End: 31},
		},
	},
	{
		input: `hassubsequence([4,1,3],[3,1])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.Comma, Start: 19, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
			{Type: lexer.CloseSquareBracket, Start: 21, End: 22},
			{Type: lexer.Comma, Start: 22, End: 23},
			{Type: lexer.OpenSquareBracket, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.CloseSquareBracket, Start: 27, End: 28},
			{Type: lexer.CloseParens, Start: 28, End: 29},
		},
	},
	{
		input: `hassubsequence([1,2],[1,1,2])`,
		expected: []lexer.Token{
			{Type: lexer.HasSubsequence, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.OpenSquareBracket, Start: 15, End: 16},
			{Type: lexer.IntegerLiteral, Start: 16, End: 17},
			{Type: lexer.Comma, Start: 17, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
			{Type: lexer.CloseSquareBracket, Start: 19, End: 20},
			{Type: lexer.Comma, Start: 20, End: 21},
			{Type: lexer.OpenSquareBracket, Start: 21, End: 22},
			{Type: lexer.IntegerLiteral, Start: 22, End: 23},
			{Type: lexer.Comma, Start: 23, End: 24},
			{Type: lexer.IntegerLiteral, Start: 24, End: 25},
			{Type: lexer.Comma, Start: 25, End: 26},
			{Type: lexer.IntegerLiteral, Start: 26, End: 27},
			{Type: lexer.CloseSquareBracket, Start: 27, End: 28},
			{Type: lexer.CloseParens, Start: 28, End: 29},
		},
	},
	{
		input: `matchesPattern(CompanyName,'%5EA.*e$')`,
		expected: []lexer.Token{
			{Type: lexer.MatchesPattern, Start: 0, End: 14},
			{Type: lexer.OpenParens, Start: 14, End: 15},
			{Type: lexer.UnquotedString, Start: 15, End: 26},
			{Type: lexer.Comma, Start: 26, End: 27},
			{Type: lexer.SingleQuotedString, Start: 27, End: 37},
			{Type: lexer.CloseParens, Start: 37, End: 38},
		},
	},
	{
		input: `tolower(CompanyName) eq 'alfreds futterkiste'`,
		expected: []lexer.Token{
			{Type: lexer.ToLower, Start: 0, End: 7},
			{Type: lexer.OpenParens, Start: 7, End: 8},
			{Type: lexer.UnquotedString, Start: 8, End: 19},
			{Type: lexer.CloseParens, Start: 19, End: 20},
			{Type: lexer.Equals, Start: 21, End: 24},
			{Type: lexer.SingleQuotedString, Start: 24, End: 45},
		},
	},
	{
		input: `toupper(CompanyName) eq 'ALFREDS FUTTERKISTE'`,
		expected: []lexer.Token{
			{Type: lexer.ToUpper, Start: 0, End: 7},
			{Type: lexer.OpenParens, Start: 7, End: 8},
			{Type: lexer.UnquotedString, Start: 8, End: 19},
			{Type: lexer.CloseParens, Start: 19, End: 20},
			{Type: lexer.Equals, Start: 21, End: 24},
			{Type: lexer.SingleQuotedString, Start: 24, End: 45},
		},
	},
	{
		input: `trim(CompanyName) eq CompanyName`,
		expected: []lexer.Token{
			{Type: lexer.Trim, Start: 0, End: 4},
			{Type: lexer.OpenParens, Start: 4, End: 5},
			{Type: lexer.UnquotedString, Start: 5, End: 16},
			{Type: lexer.CloseParens, Start: 16, End: 17},
			{Type: lexer.Equals, Start: 18, End: 21},
			{Type: lexer.UnquotedString, Start: 21, End: 32},
		},
	},
	{
		input: `day(BirthDate) eq 8`,
		expected: []lexer.Token{
			{Type: lexer.Day, Start: 0, End: 3},
			{Type: lexer.OpenParens, Start: 3, End: 4},
			{Type: lexer.UnquotedString, Start: 4, End: 13},
			{Type: lexer.CloseParens, Start: 13, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 19},
		},
	},
	{
		input: `fractionalseconds(BirthDate) lt 0.1`,
		expected: []lexer.Token{
			{Type: lexer.FractionalSeconds, Start: 0, End: 17},
			{Type: lexer.OpenParens, Start: 17, End: 18},
			{Type: lexer.UnquotedString, Start: 18, End: 27},
			{Type: lexer.CloseParens, Start: 27, End: 28},
			{Type: lexer.LessThan, Start: 29, End: 32},
			{Type: lexer.FloatingPointLiteral, Start: 32, End: 35},
		},
	},
	{
		input: `hour(BirthDate) eq 4`,
		expected: []lexer.Token{
			{Type: lexer.Hour, Start: 0, End: 4},
			{Type: lexer.OpenParens, Start: 4, End: 5},
			{Type: lexer.UnquotedString, Start: 5, End: 14},
			{Type: lexer.CloseParens, Start: 14, End: 15},
			{Type: lexer.Equals, Start: 16, End: 19},
			{Type: lexer.IntegerLiteral, Start: 19, End: 20},
		},
	},
	{
		input: `minute(BirthDate) eq 40`,
		expected: []lexer.Token{
			{Type: lexer.Minute, Start: 0, End: 6},
			{Type: lexer.OpenParens, Start: 6, End: 7},
			{Type: lexer.UnquotedString, Start: 7, End: 16},
			{Type: lexer.CloseParens, Start: 16, End: 17},
			{Type: lexer.Equals, Start: 18, End: 21},
			{Type: lexer.IntegerLiteral, Start: 21, End: 23},
		},
	},
	{
		input: `month(BirthDate) eq 5`,
		expected: []lexer.Token{
			{Type: lexer.Month, Start: 0, End: 5},
			{Type: lexer.OpenParens, Start: 5, End: 6},
			{Type: lexer.UnquotedString, Start: 6, End: 15},
			{Type: lexer.CloseParens, Start: 15, End: 16},
			{Type: lexer.Equals, Start: 17, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 21},
		},
	},
	{
		input: `second(BirthDate) eq 40`,
		expected: []lexer.Token{
			{Type: lexer.Second, Start: 0, End: 6},
			{Type: lexer.OpenParens, Start: 6, End: 7},
			{Type: lexer.UnquotedString, Start: 7, End: 16},
			{Type: lexer.CloseParens, Start: 16, End: 17},
			{Type: lexer.Equals, Start: 18, End: 21},
			{Type: lexer.IntegerLiteral, Start: 21, End: 23},
		},
	},
	{
		input: `year(BirthDate) eq 1971`,
		expected: []lexer.Token{
			{Type: lexer.Year, Start: 0, End: 4},
			{Type: lexer.OpenParens, Start: 4, End: 5},
			{Type: lexer.UnquotedString, Start: 5, End: 14},
			{Type: lexer.CloseParens, Start: 14, End: 15},
			{Type: lexer.Equals, Start: 16, End: 19},
			{Type: lexer.IntegerLiteral, Start: 19, End: 23},
		},
	},
	{
		input: `ceiling(Freight) eq 32`,
		expected: []lexer.Token{
			{Type: lexer.Ceiling, Start: 0, End: 7},
			{Type: lexer.OpenParens, Start: 7, End: 8},
			{Type: lexer.UnquotedString, Start: 8, End: 15},
			{Type: lexer.CloseParens, Start: 15, End: 16},
			{Type: lexer.Equals, Start: 17, End: 20},
			{Type: lexer.IntegerLiteral, Start: 20, End: 22},
		},
	},
	{
		input: `floor(Freight) eq 32`,
		expected: []lexer.Token{
			{Type: lexer.Floor, Start: 0, End: 5},
			{Type: lexer.OpenParens, Start: 5, End: 6},
			{Type: lexer.UnquotedString, Start: 6, End: 13},
			{Type: lexer.CloseParens, Start: 13, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 20},
		},
	},
	{
		input: `round(Freight) eq 32`,
		expected: []lexer.Token{
			{Type: lexer.Round, Start: 0, End: 5},
			{Type: lexer.OpenParens, Start: 5, End: 6},
			{Type: lexer.UnquotedString, Start: 6, End: 13},
			{Type: lexer.CloseParens, Start: 13, End: 14},
			{Type: lexer.Equals, Start: 15, End: 18},
			{Type: lexer.IntegerLiteral, Start: 18, End: 20},
		},
	},
	{
		input: `DiscontinuedDate eq null`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 16},
			{Type: lexer.Equals, Start: 17, End: 20},
			{Type: lexer.NullLiteral, Start: 20, End: 24},
		},
	},
	{
		input: `style has Sales.Pattern'Yellow'`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 5},
			{Type: lexer.Has, Start: 6, End: 10},
			{Type: lexer.UnquotedString, Start: 10, End: 23},
			{Type: lexer.SingleQuotedString, Start: 23, End: 31},
		},
	},
	{
		input: `year eq ':0'`,
		expected: []lexer.Token{
			{Type: lexer.UnquotedString, Start: 0, End: 4},
			{Type: lexer.Equals, Start: 5, End: 8},
			{Type: lexer.SingleQuotedString, Start: 8, End: 12},
		},
	},
}

func TestToken(t *testing.T) {
	t.Parallel()
	for _, test := range tests {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			l := lexer.NewLexer(tc.input)
			for _, expected := range tc.expected {
				token, err := l.NextToken()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token == nil {
					t.Fatal("unexpected nil token")
					return
				}
				if token.Type != expected.Type {
					t.Errorf("expected type %v, got %v", expected.Type, token.Type)
				}
				if token.Start != expected.Start {
					t.Errorf("expected start %v, got %v", expected.Start, token.Start)
				}
				if token.End != expected.End {
					t.Errorf("expected end %v, got %v", expected.End, token.End)
				}
			}
		})
	}
}

func BenchmarkToken(b *testing.B) {
	for _, test := range tests {
		tc := test
		b.Run(test.input, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				l := lexer.NewLexer(tc.input)
				for range tc.expected {
					_, _ = l.NextToken()
				}
			}
		})
	}
}

func FuzzToken(f *testing.F) {
	for _, test := range tests {
		f.Add(test.input)
	}
	f.Fuzz(func(_ *testing.T, input string) {
		l := lexer.NewLexer(input)
		for {
			token, err := l.NextToken()
			if err != nil {
				break
			}
			if token == nil {
				break
			}
		}
	})
}
