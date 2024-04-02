package filter

import "github.com/pboyd04/godata/filter/parser"

type Filter struct {
	myParser *parser.Parser
}

func NewFilter(input string) (*Filter, error) {
	myParser, err := parser.NewParser(input)
	if err != nil {
		return nil, err
	}
	return &Filter{myParser: myParser}, nil
}

func MustCompile(input string) *Filter {
	myParser, err := parser.NewParser(input)
	if err != nil {
		panic(err)
	}
	return &Filter{myParser: myParser}
}

func (f *Filter) GetDBQuery(language string) (interface{}, error) {
	return f.myParser.GetDBQuery(language)
}

func (f *Filter) GetDBQueryWithReplacement(language string, a ...interface{}) (interface{}, error) {
	return f.myParser.GetDBQueryWithReplacement(language, a...)
}

func (f *Filter) GetCopyWithReplacements(a ...interface{}) (*Filter, error) {
	myParser, err := f.myParser.ReplaceOperands(a...)
	if err != nil {
		return nil, err
	}
	return &Filter{myParser: myParser}, nil
}
