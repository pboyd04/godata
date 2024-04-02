package gorm

import (
	"strings"

	"github.com/pboyd04/godata/filter/lexer"
	"github.com/pboyd04/godata/filter/parser"
)

const likeStr = " LIKE ?"

type Parser struct {
}

func init() {
	// Register the parser
	parser.RegisterParser("gorm", &Parser{})
}

func (p *Parser) GetDBQuery(common *parser.Parser) (interface{}, error) {
	op, err := common.GetOperation()
	if err != nil {
		return nil, err
	}
	return p.getGormQuery(op)
}

func (p *Parser) GetDBQueryWithReplacement(common *parser.Parser, a ...interface{}) (interface{}, error) {
	op, err := common.GetOperation()
	if err != nil {
		return nil, err
	}
	op, err = op.ReplaceOperands(a...)
	if err != nil {
		return nil, err
	}
	return p.getGormQuery(op)
}

//nolint:funlen,cyclop,forcetypeassert
func (p *Parser) getGormQuery(op *parser.Operation) ([]interface{}, error) {
	operands, err := p.getGormOperands(op.Operands)
	if err != nil {
		return nil, err
	}
	//nolint:exhaustive // This won't cover everything and will use the default case to catch errors
	switch op.Operator {
	case lexer.Equals:
		return []interface{}{operands[0].(string) + " = ?", operands[1]}, nil
	case lexer.NotEquals:
		return []interface{}{operands[0].(string) + " != ?", operands[1]}, nil
	case lexer.GreaterThan:
		return []interface{}{operands[0].(string) + " > ?", operands[1]}, nil
	case lexer.GreaterThanOrEqual:
		return []interface{}{operands[0].(string) + " >= ?", operands[1]}, nil
	case lexer.LessThan:
		return []interface{}{operands[0].(string) + " < ?", operands[1]}, nil
	case lexer.LessThanOrEqual:
		return []interface{}{operands[0].(string) + " <= ?", operands[1]}, nil
	case lexer.And:
		clause1 := operands[0].([]interface{})
		clause2 := operands[1].([]interface{})
		ret := []interface{}{clause1[0].(string) + " AND " + clause2[0].(string)}
		ret = append(ret, clause1[1:]...)
		ret = append(ret, clause2[1:]...)
		return ret, nil
	case lexer.Or:
		clause1 := operands[0].([]interface{})
		clause2 := operands[1].([]interface{})
		ret := []interface{}{clause1[0].(string) + " OR " + clause2[0].(string)}
		ret = append(ret, clause1[1:]...)
		ret = append(ret, clause2[1:]...)
		return ret, nil
	case lexer.In:
		ret := []interface{}{operands[0].(string) + " IN ?"}
		ret = append(ret, operands[1:]...)
		return ret, nil
	case lexer.Contains:
		return []interface{}{operands[0].(string) + likeStr, "%" + operands[1].(string) + "%"}, nil
	case lexer.EndsWith:
		return []interface{}{operands[0].(string) + likeStr, "%" + operands[1].(string)}, nil
	case lexer.StartsWith:
		return []interface{}{operands[0].(string) + likeStr, operands[1].(string) + "%"}, nil
	case lexer.Not:
		return insertNotOp(operands[0])
	default:
		return nil, newUnsupportedOperatorError(op.Operator)
	}
}

func insertNotOp(s interface{}) ([]interface{}, error) {
	clause, ok := s.([]interface{})
	if !ok {
		return nil, newUnsupportedOperandError(s)
	}
	str, ok := clause[0].(string)
	if !ok {
		return nil, newUnsupportedOperandError(s)
	}
	parts := strings.Split(str, " ")
	if len(parts) >= 2 {
		ret := []interface{}{parts[0] + " NOT " + strings.Join(parts[1:], " ")}
		ret = append(ret, clause[1:]...)
		return ret, nil
	}
	ret := []interface{}{"NOT " + str}
	ret = append(ret, clause[1:]...)
	return ret, nil
}

func (p *Parser) getGormOperands(operands []parser.Operand) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, operand := range operands {
		tmp, err := p.getGormOperand(operand)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (p *Parser) getGormOperand(operand parser.Operand) (interface{}, error) {
	data, err := operand.GetData()
	if err != nil {
		return nil, err
	}
	switch op := data.(type) {
	case string, float64, int, map[string]interface{}:
		return op, nil
	case *parser.Operation:
		inner, err := p.getGormQuery(op)
		if err != nil {
			return nil, err
		}
		return inner, nil
	case []parser.Operand:
		tmp := make([]interface{}, 0)
		for _, o := range op {
			inner, err := p.getGormOperands([]parser.Operand{o})
			if err != nil {
				return nil, err
			}
			if len(inner) != 1 {
				tmp = append(tmp, inner)
			} else {
				tmp = append(tmp, inner[0])
			}
		}
		return tmp, nil
	default:
		return nil, newUnsupportedOperandError(op)
	}
}
