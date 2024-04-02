package gorm_test

import (
	"testing"

	"github.com/pboyd04/godata/filter/parser"
	"github.com/stretchr/testify/assert"

	_ "github.com/pboyd04/godata/filter/parser/gorm"
)

type testData struct {
	input          string
	expectedOutput []interface{}
}

//nolint:gochecknoglobals // Just test data
var testCases = []testData{
	{
		input:          "Name eq 'Milk'",
		expectedOutput: []interface{}{"Name = ?", "Milk"},
	},
	{
		input:          "(Name eq 'Milk')",
		expectedOutput: []interface{}{"Name = ?", "Milk"},
	},
	{
		input:          "Name ne 'Milk'",
		expectedOutput: []interface{}{"Name != ?", "Milk"},
	},
	{
		input:          "Name gt 'Milk'",
		expectedOutput: []interface{}{"Name > ?", "Milk"},
	},
	{
		input:          "Name ge 'Milk'",
		expectedOutput: []interface{}{"Name >= ?", "Milk"},
	},
	{
		input:          "Name lt 'Milk'",
		expectedOutput: []interface{}{"Name < ?", "Milk"},
	},
	{
		input:          "Name le 'Milk'",
		expectedOutput: []interface{}{"Name <= ?", "Milk"},
	},
	{
		input:          "Name eq 'Milk' and Price lt 2.55",
		expectedOutput: []interface{}{"Name = ? AND Price < ?", "Milk", 2.55},
	},
	{
		input:          "Name EQ 'Milk' AND Price LT 2.55",
		expectedOutput: []interface{}{"Name = ? AND Price < ?", "Milk", 2.55},
	},
	{
		input:          "Name eq 'Milk' AND Price lt 2.55",
		expectedOutput: []interface{}{"Name = ? AND Price < ?", "Milk", 2.55},
	},
	{
		input:          "Name eq 'Milk' or Price lt 2.55",
		expectedOutput: []interface{}{"Name = ? OR Price < ?", "Milk", 2.55},
	},
	{
		input:          "Name in ('Milk', 'Cheese')",
		expectedOutput: []interface{}{"Name IN ?", []interface{}{"Milk", "Cheese"}},
	},
	{
		input:          "Name in ['Milk', 'Cheese']",
		expectedOutput: []interface{}{"Name IN ?", []interface{}{"Milk", "Cheese"}},
	},
	{
		input:          "contains(Name,'red')",
		expectedOutput: []interface{}{"Name LIKE ?", "%red%"},
	},
	{
		input:          `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
		expectedOutput: []interface{}{"Address = ?", map[string]interface{}{"City": "Redmond", "State": "WA", "Street": "NE 40th", "ZipCode": "98052"}},
	},
	{
		input:          "endswith(Name,'ilk')",
		expectedOutput: []interface{}{"Name LIKE ?", "%ilk"},
	},
	{
		input:          "not endswith(Name,'ilk')",
		expectedOutput: []interface{}{"Name NOT LIKE ?", "%ilk"},
	},
	{
		input:          "startswith(CompanyName,'Futterkiste')",
		expectedOutput: []interface{}{"CompanyName LIKE ?", "Futterkiste%"},
	},
}

func TestGorm(t *testing.T) {
	t.Parallel()
	for _, test := range testCases {
		tc := test
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			common, err := parser.NewParser(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			ptrStr, err := common.GetDBQuery("gorm")
			if err != nil {
				t.Fatal(err)
			}
			res, ok := ptrStr.([]interface{})
			if !ok {
				t.Fatalf("expected []interface, got %T", ptrStr)
			}
			assert.ElementsMatch(t, tc.expectedOutput, res)
		})
	}
}
