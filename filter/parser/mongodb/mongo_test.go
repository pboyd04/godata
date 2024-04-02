package mongodb_test

import (
	"encoding/json"
	"testing"

	"github.com/pboyd04/godata/filter/parser"
	"go.mongodb.org/mongo-driver/bson"
)

type testData struct {
	input                 string
	expectedMongoJSONText string
}

type testDataReplace struct {
	input                 string
	expectedMongoJSONText string
	replacements          []interface{}
	expectedReplacedText  string
}

//nolint:gochecknoglobals // Just test data
var testCases = []testData{
	{
		input:                 "true",
		expectedMongoJSONText: "{}",
	},
	{
		input:                 "false",
		expectedMongoJSONText: `{"_id":{"$type":"string"}}`,
	},
	{
		input:                 "Name eq 'Milk'",
		expectedMongoJSONText: `{"Name":{"$eq":"Milk"}}`,
	},
	{
		input:                 "(Name eq 'Milk')",
		expectedMongoJSONText: `{"Name":{"$eq":"Milk"}}`,
	},
	{
		input:                 "Name ne 'Milk'",
		expectedMongoJSONText: `{"Name":{"$ne":"Milk"}}`,
	},
	{
		input:                 "Name gt 'Milk'",
		expectedMongoJSONText: `{"Name":{"$gt":"Milk"}}`,
	},
	{
		input:                 "Name ge 'Milk'",
		expectedMongoJSONText: `{"Name":{"$gte":"Milk"}}`,
	},
	{
		input:                 "Name lt 'Milk'",
		expectedMongoJSONText: `{"Name":{"$lt":"Milk"}}`,
	},
	{
		input:                 "Name le 'Milk'",
		expectedMongoJSONText: `{"Name":{"$lte":"Milk"}}`,
	},
	{
		input:                 "Name eq 'Milk' and Price lt 2.55",
		expectedMongoJSONText: `{"$and":[{"Name":{"$eq":"Milk"}},{"Price":{"$lt":2.55}}]}`,
	},
	{
		input:                 "Name EQ 'Milk' AND Price LT 2.55",
		expectedMongoJSONText: `{"$and":[{"Name":{"$eq":"Milk"}},{"Price":{"$lt":2.55}}]}`,
	},
	{
		input:                 "Name eq 'Milk' AND Price lt 2.55",
		expectedMongoJSONText: `{"$and":[{"Name":{"$eq":"Milk"}},{"Price":{"$lt":2.55}}]}`,
	},
	{
		input:                 "Name eq 'Milk' or Price lt 2.55",
		expectedMongoJSONText: `{"$or":[{"Name":{"$eq":"Milk"}},{"Price":{"$lt":2.55}}]}`,
	},
	{
		input:                 "Name in ('Milk', 'Cheese')",
		expectedMongoJSONText: `{"Name":{"$in":["Milk","Cheese"]}}`,
	},
	{
		input:                 "Name in ['Milk', 'Cheese']",
		expectedMongoJSONText: `{"Name":{"$in":["Milk","Cheese"]}}`,
	},
	{
		input:                 "_id eq 6206b158000e1859781d5e16",
		expectedMongoJSONText: `{"_id":{"$eq":{"$oid":"6206b158000e1859781d5e16"}}}`,
	},
	{
		input:                 "contains(Name,'red')",
		expectedMongoJSONText: `{"Name":{"$regex":"red"}}`,
	},
	{
		input:                 `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
		expectedMongoJSONText: `{"Address":{"$eq":{"City":"Redmond","State":"WA","Street":"NE 40th","ZipCode":"98052"}}}`,
	},
	{
		input:                 "endswith(Name,'ilk')",
		expectedMongoJSONText: `{"Name":{"$regex":"ilk$"}}`,
	},
	{
		input:                 "not endswith(Name,'ilk')",
		expectedMongoJSONText: `{"Name":{"$not":{"$regex":"ilk$"}}}`,
	},
	{
		input: "length(CompanyName) eq 19",
		// The numberDecimal thing is a bson-ism. Sending json to mongo CLI like this won't work, but this is correct
		expectedMongoJSONText: `{"$expr":{"$eq":[{"$strLenCP":"$CompanyName"},{"$numberDecimal":"19"}]}}`,
	},
	{
		input:                 "startswith(CompanyName,'Futterkiste')",
		expectedMongoJSONText: `{"CompanyName":{"$regex":"^Futterkiste"}}`,
	},
	{
		input:                 `hassubset(Names,["Milk", "Cheese"])`,
		expectedMongoJSONText: `{"Names":{"$all":["Milk","Cheese"]}}`,
	},
}

//nolint:gochecknoglobals // Just test data
var testCasesReplace = []testDataReplace{
	{
		input:                 "year eq ':0'",
		expectedMongoJSONText: `{"year":{"$eq":":0"}}`,
		replacements:          []interface{}{2025},
		expectedReplacedText:  `{"year":{"$eq":2025}}`,
	},
	{
		input:                 "year eq ':1' and id eq ':0'",
		expectedMongoJSONText: `{"$and":[{"year":{"$eq":":1"}},{"id":{"$eq":":0"}}]}`,
		replacements:          []interface{}{"test", 2025},
		expectedReplacedText:  `{"$and":[{"year":{"$eq":2025}},{"id":{"$eq":"test"}}]}`,
	},
}

func TestMongo(t *testing.T) {
	t.Parallel()
	for _, test := range testCases {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			parser, err := parser.NewParser(tc.input)
			if err != nil {
				t.Errorf("error creating parser: %s", err)
			}
			res, err := parser.GetDBQuery("mongodb")
			if err != nil {
				t.Errorf("error getting db query: %s", err)
			}
			mongoFilter, ok := res.(bson.D)
			if !ok {
				t.Errorf("result is not a bson.D")
			}
			bytes, _ := bson.MarshalExtJSON(mongoFilter, false, false)
			// Sort the keys...
			jsonData := make(map[string]interface{})
			_ = json.Unmarshal(bytes, &jsonData)
			bytes, _ = json.Marshal(jsonData)
			if string(bytes) != tc.expectedMongoJSONText {
				t.Fatal("Filter", tc.input, "parsed to", string(bytes), "expected", tc.expectedMongoJSONText)
			}
		})
	}
}

func TestReplacement(t *testing.T) {
	t.Parallel()
	for _, test := range testCasesReplace {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			myParser, err := parser.NewParser(tc.input)
			if err != nil {
				t.Errorf("error creating parser: %s", err)
			}
			noReplaceRes, err := myParser.GetDBQuery("mongodb")
			if err != nil {
				t.Fatalf("error getting db query: %s", err)
			}
			mongoFilter, ok := noReplaceRes.(bson.D)
			if !ok {
				t.Fatalf("result is not a bson.D")
			}
			bytes, _ := bson.MarshalExtJSON(mongoFilter, false, false)
			if string(bytes) != tc.expectedMongoJSONText {
				t.Error("Filter without replacement", tc.input, "parsed to", string(bytes), "expected", tc.expectedMongoJSONText)
			}
			res, err := myParser.GetDBQueryWithReplacement("mongodb", tc.replacements...)
			if err != nil {
				t.Errorf("error getting db query: %s", err)
			}
			mongoFilter, ok = res.(bson.D)
			if !ok {
				t.Errorf("result is not a bson.D")
			}
			bytes, _ = bson.MarshalExtJSON(mongoFilter, false, false)
			// Sort the keys...
			jsonData := make(map[string]interface{})
			_ = json.Unmarshal(bytes, &jsonData)
			bytes, _ = json.Marshal(jsonData)
			if string(bytes) != tc.expectedReplacedText {
				t.Fatal("Filter", tc.input, "parsed to", string(bytes), "expected", tc.expectedReplacedText)
			}
		})
	}
}

func BenchmarkMongo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			tc := test
			b.Run(test.input, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					parser, err := parser.NewParser(tc.input)
					if err != nil {
						b.Errorf("error creating parser: %s", err)
					}
					_, err = parser.GetDBQuery("mongodb")
					if err != nil {
						b.Errorf("error getting db query: %s", err)
					}
				}
			})
		}
	}
}

func BenchmarkReplacement(b *testing.B) {
	// Setup the replacements outside the benchmark loop, that's the point of the function
	replacements := make([]*parser.Parser, 0)
	for _, test := range testCasesReplace {
		myParser, err := parser.NewParser(test.input)
		if err != nil {
			b.Fatalf("error creating parser: %s", err)
		}
		replacements = append(replacements, myParser)
	}
	for i, test := range testCasesReplace {
		tc := test
		index := i
		b.Run(test.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := replacements[index].GetDBQueryWithReplacement("mongodb", tc.replacements...)
				if err != nil {
					b.Errorf("error getting db query: %s", err)
				}
			}
		})
	}
}
