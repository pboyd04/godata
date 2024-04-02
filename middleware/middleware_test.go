package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/pboyd04/godata/filter/parser/mysql"
	"github.com/pboyd04/godata/middleware"
	"github.com/pboyd04/godata/orderby"
)

type testData struct {
	input           string
	expectedFilter  string
	expectedSelect  *[]string
	expectedOrderBy *orderby.OrderBy
	expectedTop     int64
	expectedSkip    int64
	expectedCount   bool
}

//nolint:gochecknoglobals // Just test data
var tests = []testData{
	{
		input:           "$filter=Name%20eq%20'Bob'",
		expectedFilter:  "`Name`='Bob'",
		expectedSelect:  nil,
		expectedOrderBy: nil,
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$select=Name,Test",
		expectedFilter:  "",
		expectedSelect:  &[]string{"Name", "Test"},
		expectedOrderBy: nil,
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$orderby=Name",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: &orderby.OrderBy{OrderItem: []orderby.OrderItem{{Property: "Name", Direction: orderby.ASC}}},
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$orderby=Name%20DESC",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: &orderby.OrderBy{OrderItem: []orderby.OrderItem{{Property: "Name", Direction: orderby.DESC}}},
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$top=10",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: nil,
		expectedTop:     10,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$skip=100",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: nil,
		expectedTop:     -1,
		expectedSkip:    100,
		expectedCount:   false,
	},
	{
		input:           "$count=true",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: nil,
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   true,
	},
	{
		input:           "$count=false",
		expectedFilter:  "",
		expectedSelect:  nil,
		expectedOrderBy: nil,
		expectedTop:     -1,
		expectedSkip:    -1,
		expectedCount:   false,
	},
	{
		input:           "$filter=Name%20eq%20'Bob'&$select=Name,Test&$orderby=Name&$top=10&$skip=100&$count=true",
		expectedFilter:  "`Name`='Bob'",
		expectedSelect:  &[]string{"Name", "Test"},
		expectedOrderBy: &orderby.OrderBy{OrderItem: []orderby.OrderItem{{Property: "Name", Direction: orderby.ASC}}},
		expectedTop:     10,
		expectedSkip:    100,
		expectedCount:   true,
	},
}

func TestMiddlewareEmpty(t *testing.T) {
	t.Parallel()
	middleware := middleware.NewOdataMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		odata := middleware.GetOdataFromContext(r.Context())
		if odata == nil {
			t.Fatal("Missing odata filter")
		}
		if odata.Filter != nil {
			t.Error("Filter should be nil")
		}
		if odata.Select != nil {
			t.Error("Select should be nil")
		}
		if odata.OrderBy != nil {
			t.Error("OrderBy should be nil")
		}
		if odata.Top != -1 {
			t.Error("Top should be -1")
		}
		if odata.Skip != -1 {
			t.Error("Skip should be -1")
		}
		if odata.Count != false {
			t.Error("Count should be false")
		}
	}))
	middleware.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", nil))
}

//nolint:gocognit,cyclop // This is just a test
func TestMiddleware(t *testing.T) {
	t.Parallel()
	for _, test := range tests {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			middleware := middleware.NewOdataMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				odata := middleware.GetOdataFromContext(r.Context())
				if odata == nil {
					t.Fatal("Missing odata filter")
				}
				if odata.Filter != nil && tc.expectedFilter == "" {
					t.Error("Filter should be nil")
				} else if odata.Filter != nil {
					data, err := odata.Filter.GetDBQuery("mysql")
					if err != nil {
						t.Fatal(err)
					}
					//nolint:forcetypeassert // Just test code
					if data.(string) != tc.expectedFilter {
						//nolint:forcetypeassert // Just test code
						t.Errorf("Filter should be %s got %s", tc.expectedFilter, data.(string))
					}
				}
				if !sliceEq(odata.Select, tc.expectedSelect) {
					t.Errorf("Select should be %v got %v", tc.expectedSelect, odata.Select)
				}
				if !orderByEq(odata.OrderBy, tc.expectedOrderBy) {
					t.Errorf("OrderBy should be %v got %v", tc.expectedOrderBy, odata.OrderBy)
				}
				if odata.Top != tc.expectedTop {
					t.Errorf("Top should be %d", tc.expectedTop)
				}
				if odata.Skip != tc.expectedSkip {
					t.Errorf("Skip should be %d", tc.expectedSkip)
				}
				if odata.Count != tc.expectedCount {
					t.Errorf("Count should be %t", tc.expectedCount)
				}
			}))
			req := httptest.NewRequest(http.MethodGet, "/test?"+tc.input, nil)
			middleware.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

func BenchmarkMiddlewareServeHTTP(b *testing.B) {
	for _, test := range tests {
		tc := test
		b.Run(test.input, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			req := httptest.NewRequest(http.MethodGet, "/test?"+tc.input, nil)
			middleware := middleware.NewOdataMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
			for i := 0; i < b.N; i++ {
				middleware.ServeHTTP(httptest.NewRecorder(), req)
			}
		})
	}
}

func FuzzMiddleware(f *testing.F) {
	for _, test := range tests {
		f.Add(test.input)
	}
	f.Fuzz(func(_ *testing.T, input string) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.URL.RawQuery = input
		middleware := middleware.NewOdataMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
		middleware.ServeHTTP(httptest.NewRecorder(), req)
	})
}

func sliceEq(a, b *[]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	for i, v := range *a {
		if v != (*b)[i] {
			return false
		}
	}
	return true
}

func orderByEq(a, b *orderby.OrderBy) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a.OrderItem) != len(b.OrderItem) {
		return false
	}
	for i, v := range a.OrderItem {
		if v != b.OrderItem[i] {
			return false
		}
	}
	return true
}
