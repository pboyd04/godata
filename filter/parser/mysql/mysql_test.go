package mysql_test

import (
	"testing"

	"github.com/pboyd04/godata/filter/parser"
)

type testData struct {
	input           string
	expectedSQLText string
}

//nolint:gochecknoglobals // Just test data
var testCases = []testData{
	{
		input:           "true",
		expectedSQLText: `1=1`,
	},
	{
		input:           "false",
		expectedSQLText: `1=0`,
	},
	{
		input:           "Name eq 'Milk'",
		expectedSQLText: "`Name`='Milk'",
	},
	{
		input:           "(Name eq 'Milk')",
		expectedSQLText: "`Name`='Milk'",
	},
	{
		input:           "Name ne 'Milk'",
		expectedSQLText: "`Name`!='Milk'",
	},
	{
		input:           "Name gt 'Milk'",
		expectedSQLText: "`Name`>'Milk'",
	},
	{
		input:           "Name ge 'Milk'",
		expectedSQLText: "`Name`>='Milk'",
	},
	{
		input:           "Name lt 'Milk'",
		expectedSQLText: "`Name`<'Milk'",
	},
	{
		input:           "Name le 'Milk'",
		expectedSQLText: "`Name`<='Milk'",
	},
	{
		input:           "Name eq 'Milk' and Price lt 2.55",
		expectedSQLText: "`Name`='Milk' AND `Price`<2.55",
	},
	{
		input:           "Name EQ 'Milk' AND Price LT 2.55",
		expectedSQLText: "`Name`='Milk' AND `Price`<2.55",
	},
	{
		input:           "Name eq 'Milk' AND Price lt 2.55",
		expectedSQLText: "`Name`='Milk' AND `Price`<2.55",
	},
	{
		input:           "Name eq 'Milk' or Price lt 2.55",
		expectedSQLText: "`Name`='Milk' OR `Price`<2.55",
	},
	{
		input:           "Name in ('Milk', 'Cheese')",
		expectedSQLText: "`Name` IN ('Milk','Cheese')",
	},
	{
		input:           "Name in ['Milk', 'Cheese']",
		expectedSQLText: "`Name` IN ('Milk','Cheese')",
	},
	{
		input:           "contains(Name,'red')",
		expectedSQLText: "`Name` LIKE '%red%'",
	},
	{
		input:           `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
		expectedSQLText: "`Address`='{\\\"City\\\":\\\"Redmond\\\",\\\"State\\\":\\\"WA\\\",\\\"Street\\\":\\\"NE 40th\\\",\\\"ZipCode\\\":\\\"98052\\\"}'",
	},
	{
		input:           "endswith(Name,'ilk')",
		expectedSQLText: "`Name` LIKE '%ilk'",
	},
	{
		input:           "not endswith(Name,'ilk')",
		expectedSQLText: "`Name` NOT LIKE '%ilk'",
	},
	{
		input:           "length(CompanyName) eq 19",
		expectedSQLText: "LENGTH(`CompanyName`)=19",
	},
	{
		input:           "startswith(CompanyName,'Futterkiste')",
		expectedSQLText: "`CompanyName` LIKE 'Futterkiste%'",
	},
	{
		input:           `hassubset(Names,["Milk", "Cheese"])`,
		expectedSQLText: "JSON_CONTAINS(`Names`,'[\"Milk\",\"Cheese\"]')", // This is mysql syntax. Don't copy to other SQL parsers
	},
	{
		input:           `Price add 2.45 eq 5.00`,
		expectedSQLText: "`Price`+2.45=5",
	},
	{
		input:           `Price sub 0.55 eq 2.00`,
		expectedSQLText: "`Price`-0.55=2",
	},
	{
		input:           `Price mul 2.0 eq 5.10`,
		expectedSQLText: "`Price`*2=5.1",
	},
	{
		input:           `Price div 2.55 eq 1`,
		expectedSQLText: "`Price`/2.55=1",
	},
	{
		input:           `Rating div 2 eq 2`,
		expectedSQLText: "`Rating` DIV 2=2",
	},
	{
		input:           `Rating divby 2 eq 2.5`,
		expectedSQLText: "`Rating`/2=2.5",
	},
	{
		input:           `Rating mod 5 eq 0`,
		expectedSQLText: "`Rating` MOD 5=0",
	},
}

func TestMySQL(t *testing.T) {
	t.Parallel()
	for _, test := range testCases {
		tc := test
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			common, err := parser.NewParser(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			ptrStr, err := common.GetDBQuery("mysql")
			if err != nil {
				t.Fatal(err)
			}
			str, ok := ptrStr.(string)
			if !ok {
				t.Fatalf("expected string, got %T", ptrStr)
			}
			if str != tc.expectedSQLText {
				t.Errorf("expected %q, got %q", tc.expectedSQLText, str)
			}
		})
	}
}

func BenchmarkMySQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			tc := test
			b.Run(test.input, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					common, err := parser.NewParser(tc.input)
					if err != nil {
						b.Fatal(err)
					}
					_, err = common.GetDBQuery("mysql")
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
