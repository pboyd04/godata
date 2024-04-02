package parser

import (
	"encoding/json"

	"github.com/pboyd04/godata/filter/lexer"
)

type Parser struct {
	lexer  *lexer.Lexer
	tokens []*lexer.Token
	op     *Operation
}

type IParser interface {
	GetDBQuery(p *Parser) (interface{}, error)
	GetDBQueryWithReplacement(p *Parser, a ...interface{}) (interface{}, error)
}

//nolint:gochecknoglobals // This is a map of parsers, required to let other parsers register themselves
var parsers = map[string]IParser{}

func NewParser(input string) (*Parser, error) {
	myLexer := lexer.NewLexer(input)
	ret := &Parser{lexer: myLexer, tokens: []*lexer.Token{}}
	token, err := myLexer.NextToken()
	if err != nil {
		return nil, err
	}
	for token != nil {
		ret.tokens = append(ret.tokens, token)
		token, err = myLexer.NextToken()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

type Operator int

const (
	NoOp Operator = -1
)

type Operand interface {
	GetData() (interface{}, error)
}

type SliceOperand struct {
	Slice []Operand
}

type ObjectOperand struct {
	Properties string
}

type Operation struct {
	Operator Operator  `json:"operator"`
	Operands []Operand `json:"operands,omitempty"`
}

type tokenComparer func(*lexer.Token) bool

func (p *Parser) GetOperation() (*Operation, error) {
	if p.op != nil {
		return p.op, nil
	}
	// Make all the tokens one big group...
	tokenGroup := newTokenGroup(p.tokens)
	// Handle parentheses
	tokenGroup.parentheses()
	// Group {}
	err := tokenGroup.handleObjects()
	if err != nil {
		return nil, err
	}
	// Start figuring out ops...
	err = tokenGroup.processTokenGroupOps()
	if err != nil {
		return nil, err
	}
	return tokenGroup.flatten(nil)
}

func (o *Operation) flatten() error {
	for i := 0; i < len(o.Operands); i++ {
		switch operand := o.Operands[i].(type) {
		case *tokenGroup:
			if o.Operator.hasParameters() {
				// This is a method call
				o.Operands = append(o.Operands, operand.children...)
				// Remove all token groups from the parent
				for i := 0; i < len(o.Operands); i++ {
					_, ok := o.Operands[i].(*tokenGroup)
					if ok {
						o.Operands = append(o.Operands[:i], o.Operands[i+1:]...)
					}
				}
				return o.flatten()
			}
			op, err := operand.flatten(o)
			if err != nil {
				return operand.handleFlattenError(o, err)
			}
			o.Operands[i] = op
		case *Operation:
			err := operand.flatten()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Operation) unary() error {
	for i := 0; i < len(o.Operands); i++ {
		switch operand := o.Operands[i].(type) {
		case *tokenGroup:
			err := operand.unary()
			if err != nil {
				return err
			}
		case *Operation:
			err := operand.unary()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Operation) methodCalls() error {
	for i := 0; i < len(o.Operands); i++ {
		switch operand := o.Operands[i].(type) {
		case *tokenGroup:
			err := operand.methodCalls()
			if err != nil {
				return err
			}
		case *Operation:
			err := operand.methodCalls()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func multiplicativeCompare(token *lexer.Token) bool {
	return token.IsMultiplicative()
}

func additiveCompare(token *lexer.Token) bool {
	return token.IsAdditive()
}

func relationalCompare(token *lexer.Token) bool {
	return token.IsRelational()
}

func equalityCompare(token *lexer.Token) bool {
	return token.IsEquality()
}

func conjunctionCompare(token *lexer.Token) bool {
	return token.Type == lexer.And || token.Type == lexer.Or
}

func (o *Operation) beforeAndAfterOpProcess(compareFn tokenComparer) error {
	for i := 0; i < len(o.Operands); i++ {
		switch operand := o.Operands[i].(type) {
		case *tokenGroup:
			err := operand.beforeAndAfterOpProcess(compareFn)
			if err != nil {
				return err
			}
		case *Operation:
			err := operand.beforeAndAfterOpProcess(compareFn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Operation) GetData() (interface{}, error) {
	return o, nil
}

func (s *SliceOperand) GetData() (interface{}, error) {
	return s.Slice, nil
}

func (o *ObjectOperand) GetData() (interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(o.Properties), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (o *Operator) hasParameters() bool {
	return lexer.TokenKey(*o).HasParameters()
}

func (o *Operation) ReplaceOperands(a ...interface{}) (*Operation, error) {
	newOp, err := o.deepClone()
	if err != nil {
		return nil, err
	}
	for i, operand := range a {
		err := newOp.replaceOperand(i, operand)
		if err != nil {
			return nil, err
		}
	}
	return newOp, nil
}

func (o *Operation) deepClone() (*Operation, error) {
	newOp := &Operation{Operator: o.Operator, Operands: make([]Operand, len(o.Operands))}
	for i, operand := range o.Operands {
		switch op := operand.(type) {
		case *lexer.Token:
			newOp.Operands[i] = &lexer.Token{Type: op.Type, Text: op.Text}
		case *Operation:
			newChild, err := op.deepClone()
			if err != nil {
				return nil, err
			}
			newOp.Operands[i] = newChild
		default:
			return nil, newParserError("unknown type: %T", operand)
		}
	}
	return newOp, nil
}

func (o *Operation) replaceOperand(i int, operand interface{}) error {
	for _, inOperand := range o.Operands {
		switch op := inOperand.(type) {
		case *lexer.Token:
			if op.IsCorrectReplacement(i) {
				err := op.Replace(operand)
				if err != nil {
					return err
				}
			}
		case *Operation:
			err := op.replaceOperand(i, operand)
			if err != nil {
				return err
			}
		default:
			return newParserError("unknown type: %T", inOperand)
		}
	}
	return nil
}

func newSliceOperand(tokens []Operand) *SliceOperand {
	ret := &SliceOperand{Slice: make([]Operand, 0)}
	for _, token := range tokens {
		token, ok := token.(*lexer.Token)
		if ok && token.Type == lexer.Comma {
			// Skip commas
			continue
		}
		ret.Slice = append(ret.Slice, token)
	}
	return ret
}

func (p *Parser) GetDBQuery(language string) (interface{}, error) {
	parser, ok := parsers[language]
	if !ok {
		return nil, ErrNoSuchLanguage
	}
	return parser.GetDBQuery(p)
}

func (p *Parser) GetDBQueryWithReplacement(language string, a ...interface{}) (interface{}, error) {
	parser, ok := parsers[language]
	if !ok {
		return nil, ErrNoSuchLanguage
	}
	return parser.GetDBQueryWithReplacement(p, a...)
}

func (p *Parser) ReplaceOperands(a ...interface{}) (*Parser, error) {
	op, err := p.GetOperation()
	if err != nil {
		return nil, err
	}
	op, err = op.ReplaceOperands(a...)
	if err != nil {
		return nil, err
	}
	return &Parser{lexer: p.lexer, tokens: nil, op: op}, nil
}

func RegisterParser(name string, parser IParser) {
	parsers[name] = parser
}

func (o *Operator) MarshalJSON() ([]byte, error) {
	return json.Marshal(lexer.TokenKey(*o).String())
}
