package filter_test

import (
	"testing"

	"github.com/pboyd04/godata/filter"
)

type filterTestData struct {
	filterText string
}

//nolint:gochecknoglobals // This is just test data.
var testData = []filterTestData{
	{
		filterText: "true",
	},
	{
		filterText: "false",
	},
	{
		filterText: "Name eq 'Milk'",
	},
	{
		filterText: "(Name eq 'Milk')",
	},
	{
		filterText: "Name ne 'Milk'",
	},
	{
		filterText: "Name gt 'Milk'",
	},
	{
		filterText: "Name ge 'Milk'",
	},
	{
		filterText: "Name lt 'Milk'",
	},
	{
		filterText: "Name le 'Milk'",
	},
	{
		filterText: "Name eq 'Milk' and Price lt 2.55",
	},
	{
		filterText: "Name EQ 'Milk' AND Price LT 2.55",
	},
	{
		filterText: "Name eq 'Milk' AND Price lt 2.55",
	},
	{
		filterText: "Name eq 'Milk' or Price lt 2.55",
	},
	{
		filterText: "Name in ('Milk', 'Cheese')",
	},
	{
		filterText: "Name in ['Milk', 'Cheese']",
	},
	{
		filterText: "_id eq 6206b158000e1859781d5e16",
	},
	{
		filterText: "contains(Name,'red')",
	},
	{
		filterText: `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
	},
	{
		filterText: "endswith(Name,'ilk')",
	},
	{
		filterText: "not endswith(Name,'ilk')",
	},
	{
		filterText: "length(CompanyName) eq 19",
	},
	{
		filterText: "startswith(CompanyName,'Futterkiste')",
	},
	{
		filterText: `hassubset(Names,["Milk", "Cheese"])`,
	},
	{
		filterText: `Price add 2.45 eq 5.00`,
	},
	{
		filterText: `Price sub 0.55 eq 2.00`,
	},
	{
		filterText: `Price mul 2.0 eq 5.10`,
	},
	{
		filterText: `Price div 2.55 eq 1`,
	},
	{
		filterText: `Rating div 2 eq 2`,
	},
	{
		filterText: `Rating divby 2 eq 2.5`,
	},
	{
		filterText: `Rating mod 5 eq 0`,
	},
}

func TestFilterParsingGood(t *testing.T) {
	t.Parallel()
	for _, test := range testData {
		_, err := filter.NewFilter(test.filterText)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCaseInsensitiveAndOr(t *testing.T) {
	t.Parallel()
	_, err := filter.NewFilter("Name eq 'Milk' and Price lt 2.55")
	if err != nil {
		t.Fatal(err)
	}
	_, err = filter.NewFilter("Name eq 'Milk' AND Price lt 2.55")
	if err != nil {
		t.Fatal(err)
	}
	_, err = filter.NewFilter("Name eq 'Milk' or Price lt 2.55")
	if err != nil {
		t.Fatal(err)
	}
	_, err = filter.NewFilter("Name eq 'Milk' OR Price lt 2.55")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMustCompile(t *testing.T) {
	t.Parallel()
	for _, test := range testData {
		filter.MustCompile(test.filterText)
		// Didn't panic so all is good...
	}
}
