package orderby_test

import (
	"testing"

	"github.com/pboyd04/godata/orderby"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	input    string
	expected []orderby.OrderItem
}

//nolint:gochecknoglobals // Just test data
var tests = []testData{
	{
		input:    "",
		expected: []orderby.OrderItem{},
	},
	{
		input:    "BaseRate asc",
		expected: []orderby.OrderItem{{Property: "BaseRate", Direction: orderby.ASC}},
	},
	{
		input:    "Rating desc,BaseRate",
		expected: []orderby.OrderItem{{Property: "Rating", Direction: orderby.DESC}, {Property: "BaseRate", Direction: orderby.ASC}},
	},
}

func TestOrderBy(t *testing.T) {
	t.Parallel()
	for _, test := range tests {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			res, err := orderby.NewOrderBy(tc.input)
			if err != nil {
				t.Errorf("NewOrderBy(%s) returned an error: %v", tc.input, err)
			}
			assert.ElementsMatch(t, tc.expected, res.OrderItem)
		})
	}
}
