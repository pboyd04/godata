package golang_test

import (
	"testing"
	"time"

	"github.com/pboyd04/godata/filter/parser"
	"github.com/pboyd04/godata/filter/parser/golang"

	"github.com/stretchr/testify/assert"
)

//nolint:tagliatelle // Test data
type testInputStruct struct {
	Name      string
	JSONInput string `json:"jsonInput"`
	Int       int
	Float64   float64 `json:"Price"`
	Array     []string
	L         string `json:"City"`
	C         string `json:"Country"`
	IntArray  []int
	Date      time.Time
	TestPtr   *testInputStruct
}

//nolint:gochecknoglobals // Just test data
var testInputData = []interface{}{
	testInputStruct{
		Name:      "structuredTest",
		JSONInput: "jsonTest",
		Int:       1,
		Array:     []string{"1", "2", "3", "5"},
		IntArray:  []int{1, 2, 3, 5},
		Date:      time.Date(2022, 9, 8, 4, 0, 0, 0, time.UTC),
		TestPtr:   new(testInputStruct),
	},
	testInputStruct{
		Name:     "bob ",
		Int:      -1,
		Float64:  2.55,
		Array:    []string{"5", "2", "3", "1"},
		IntArray: []int{5, 2, 3, 1},
		Date:     time.Date(2021, 8, 9, 0, 40, 40, 10000, time.UTC),
	},
	testInputStruct{
		Name:    "Milk",
		Int:     0,
		Float64: 2.55,
		Array:   []string{"Milk", "Cheese"},
	},
	testInputStruct{
		Name:    "Milk",
		Int:     5,
		Float64: 1.1,
		Array:   []string{"Milk", "Bob"},
		L:       "Berlin",
		C:       "United States",
	},
	testInputStruct{
		Name:    "Cheese",
		Int:     4,
		Float64: 10.1,
		L:       "Berlin",
		C:       "Germany",
	},
}

type testData struct {
	input          string
	expectedOutput []interface{}
}

//nolint:gochecknoglobals // Just test data
var testCases = []testData{
	{
		input:          "true",
		expectedOutput: testInputData,
	},
	{
		input:          "false",
		expectedOutput: []interface{}{},
	},
	{
		input:          "Name eq 'Milk'",
		expectedOutput: []interface{}{testInputData[2], testInputData[3]},
	},
	{
		input:          "(Name eq 'Milk')",
		expectedOutput: []interface{}{testInputData[2], testInputData[3]},
	},
	{
		input:          "Name ne 'Milk'",
		expectedOutput: []interface{}{testInputData[0], testInputData[1], testInputData[4]},
	},
	{
		input:          "Name gt 'Milk'",
		expectedOutput: []interface{}{testInputData[0], testInputData[1]},
	},
	{
		input:          "Name ge 'Milk'",
		expectedOutput: []interface{}{testInputData[0], testInputData[1], testInputData[2], testInputData[3]},
	},
	{
		input:          "Name lt 'Milk'",
		expectedOutput: []interface{}{testInputData[4]},
	},
	{
		input:          "Name le 'Milk'",
		expectedOutput: []interface{}{testInputData[2], testInputData[3], testInputData[4]},
	},
	{
		input:          "Name eq 'Milk' and Price lt 2.55",
		expectedOutput: []interface{}{testInputData[3]},
	},
	{
		input:          "Name EQ 'Milk' AND Price LT 2.55",
		expectedOutput: []interface{}{testInputData[3]},
	},
	{
		input:          "Name eq 'Milk' AND Price lt 2.55",
		expectedOutput: []interface{}{testInputData[3]},
	},
	{
		input:          "Name eq 'Milk' AND Price eq 2.55",
		expectedOutput: []interface{}{testInputData[2]},
	},
	{
		input:          "Name eq 'Milk' or Price lt 2.55",
		expectedOutput: []interface{}{testInputData[0], testInputData[2], testInputData[3]},
	},
	{
		input:          "Name in ('Milk', 'Cheese')",
		expectedOutput: []interface{}{testInputData[2], testInputData[3], testInputData[4]},
	},
	{
		input:          "Name in ['Milk', 'Cheese']",
		expectedOutput: []interface{}{testInputData[2], testInputData[3], testInputData[4]},
	},
	{
		input:          "contains(Name,'red')",
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          "endswith(Name,'ilk')",
		expectedOutput: []interface{}{testInputData[2], testInputData[3]},
	},
	{
		input:          "not endswith(Name,'ilk')",
		expectedOutput: []interface{}{testInputData[0], testInputData[1], testInputData[4]},
	},
	{
		input:          "length(Name) gt 4",
		expectedOutput: []interface{}{testInputData[0], testInputData[4]},
	},
	{
		input:          "startswith(Name,'str')",
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `hassubset(Array,["Milk", "Cheese"])`,
		expectedOutput: []interface{}{testInputData[2]},
	},
	{
		input:          `Price add 2.45 eq 5.00`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `Price sub 0.55 eq 2.00`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `Price mul 2.0 eq 5.10`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `Price div 2.55 eq 1`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `Int div 2 eq 2`,
		expectedOutput: []interface{}{testInputData[3], testInputData[4]},
	},
	{
		input:          `Int divby 2 eq 2.5`,
		expectedOutput: []interface{}{testInputData[3]},
	},
	{
		input:          `Int mod 5 eq 0`,
		expectedOutput: []interface{}{testInputData[2], testInputData[3]},
	},
	{
		input:          `concat(concat(City,', '),Country) eq 'Berlin, Germany'`,
		expectedOutput: []interface{}{testInputData[4]},
	},
	{
		input:          `indexof(Name,'Test') eq 10`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `substring(Name,1) eq 'ob '`,
		expectedOutput: []interface{}{testInputData[1]},
	},
	{
		input:          `substring(Name,1,3) eq 'hee'`,
		expectedOutput: []interface{}{testInputData[4]},
	},
	{
		input:          `hassubsequence(Array,['1','2','3'])`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `hassubsequence(IntArray,[1,2,3])`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `matchesPattern(Name,'^[A-Z]')`,
		expectedOutput: []interface{}{testInputData[2], testInputData[3], testInputData[4]},
	},
	{
		input:          `tolower(Name) eq 'structuredtest'`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `toupper(Name) eq 'BOB '`,
		expectedOutput: []interface{}{testInputData[1]},
	},
	{
		input:          `trim(Name) eq Name`,
		expectedOutput: []interface{}{testInputData[0], testInputData[2], testInputData[3], testInputData[4]},
	},
	{
		input:          `day(Date) eq 8`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `fractionalseconds(Date) ge 0.01`,
		expectedOutput: []interface{}{testInputData[1]},
	},
	{
		input:          `hour(Date) eq 4`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `minute(Date) eq 40`,
		expectedOutput: []interface{}{testInputData[1]},
	},
	{
		input:          `month(Date) eq 9`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `second(Date) eq 40`,
		expectedOutput: []interface{}{testInputData[1]},
	},
	{
		input:          `year(Date) eq 2022`,
		expectedOutput: []interface{}{testInputData[0]},
	},
	{
		input:          `ceiling(Price) eq 3`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `floor(Price) eq 2`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `round(Price) eq 3`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2]},
	},
	{
		input:          `TestPtr eq null`,
		expectedOutput: []interface{}{testInputData[1], testInputData[2], testInputData[3], testInputData[4]},
	},
}

func TestGoLang(t *testing.T) {
	t.Parallel()
	for _, test := range testCases {
		tc := test
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			common, err := parser.NewParser(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			ptr, err := common.GetDBQuery("golang")
			if err != nil {
				t.Fatal(err)
			}
			eval, ok := ptr.(*golang.Evaluator)
			if !ok {
				t.Fatalf("expected Evaluator, got %T", ptr)
			}
			res, err := eval.FilterSlice(testInputData)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, tc.expectedOutput, res)
		})
	}
}
