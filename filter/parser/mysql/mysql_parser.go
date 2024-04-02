package mysql

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/pboyd04/godata/filter/lexer"
	"github.com/pboyd04/godata/filter/parser"
)

const badString = "ERR! NOT A STRING"

type Parser struct {
	functionMatch       *regexp.Regexp
	alreadyEscapedMatch *regexp.Regexp
}

func init() {
	// Register the parser
	parser.RegisterParser("mysql", &Parser{
		functionMatch:       regexp.MustCompile(`[A-Z]+[(]`),
		alreadyEscapedMatch: regexp.MustCompile(`\x60(\w)+\x60`),
	})
}

func (p *Parser) GetDBQuery(common *parser.Parser) (interface{}, error) {
	op, err := common.GetOperation()
	if err != nil {
		return nil, err
	}
	return p.getMySQLQuery(op)
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
	return p.getMySQLQuery(op)
}

//nolint:funlen,cyclop
func (p *Parser) getMySQLQuery(op *parser.Operation) (string, error) {
	operands, err := p.getMySQLOperands(op.Operands)
	if err != nil {
		return "", err
	}
	//nolint:exhaustive // This won't cover everything and will use the default case to catch errors
	switch op.Operator {
	case lexer.TokenTrue:
		return "1=1", nil
	case lexer.TokenFalse:
		return "1=0", nil
	case lexer.Equals:
		return p.escapeColName(operands[0]) + "=" + escapeValue(operands[1]), nil
	case lexer.NotEquals:
		return p.escapeColName(operands[0]) + "!=" + escapeValue(operands[1]), nil
	case lexer.GreaterThan:
		return p.escapeColName(operands[0]) + ">" + escapeValue(operands[1]), nil
	case lexer.GreaterThanOrEqual:
		return p.escapeColName(operands[0]) + ">=" + escapeValue(operands[1]), nil
	case lexer.LessThan:
		return p.escapeColName(operands[0]) + "<" + escapeValue(operands[1]), nil
	case lexer.LessThanOrEqual:
		return p.escapeColName(operands[0]) + "<=" + escapeValue(operands[1]), nil
	case lexer.In:
		return p.escapeColName(operands[0]) + " IN " + escapeValue(operands[1]), nil
	case lexer.And, lexer.Or:
		return p.doCombination(op.Operator, operands[0], operands[1])
	case lexer.StartsWith:
		return p.doRegex("", "%", operands[0], operands[1])
	case lexer.EndsWith:
		return p.doRegex("%", "", operands[0], operands[1])
	case lexer.Contains:
		return p.doRegex("%", "%", operands[0], operands[1])
	case lexer.Not:
		return insertNotOp(operands[0])
	case lexer.Length:
		return "LENGTH(" + p.escapeColName(operands[0]) + ")", nil
	case lexer.HasSubset:
		return "JSON_CONTAINS(" + p.escapeColName(operands[0]) + ",'" + escapeJSONValue(operands[1]) + "')", nil
	case lexer.Add:
		return p.escapeColName(operands[0]) + "+" + escapeValue(operands[1]), nil
	case lexer.Subtract:
		return p.escapeColName(operands[0]) + "-" + escapeValue(operands[1]), nil
	case lexer.Multiply:
		return p.escapeColName(operands[0]) + "*" + escapeValue(operands[1]), nil
	case lexer.Divide:
		_, ok := operands[1].(int)
		if ok {
			// This is an integer, so I need to use the DIV operator per odata spec
			return p.escapeColName(operands[0]) + " DIV " + escapeValue(operands[1]), nil
		}
		return p.escapeColName(operands[0]) + "/" + escapeValue(operands[1]), nil
	case lexer.DivideFloat:
		return p.escapeColName(operands[0]) + "/" + escapeValue(operands[1]), nil
	case lexer.Modulo:
		return p.escapeColName(operands[0]) + " MOD " + escapeValue(operands[1]), nil
	default:
		return "", newUnsupportedOperatorError(op.Operator)
	}
}

func (p *Parser) doCombination(op parser.Operator, operand0, operand1 interface{}) (string, error) {
	comb := " AND "
	if op == lexer.Or {
		comb = " OR "
	}
	strOp0, ok := operand0.(string)
	if !ok {
		return "", newParserError("attempting to combine a non-string op0")
	}
	strOp1, ok := operand1.(string)
	if !ok {
		return "", newParserError("attempting to combine a non-string op1")
	}
	return strOp0 + comb + strOp1, nil
}

func (p *Parser) doRegex(prefix, postfix string, operand0, operand1 interface{}) (string, error) {
	strOp1, ok := operand1.(string)
	if !ok {
		return "", newParserError("attempting to do a regex with a non-string value")
	}
	return p.escapeColName(operand0) + " LIKE '" + prefix + strOp1 + postfix + "'", nil
}

func (p *Parser) getMySQLOperands(operands []parser.Operand) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, operand := range operands {
		tmp, err := p.getMySQLOperand(operand)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (p *Parser) getMySQLOperand(operand parser.Operand) (interface{}, error) {
	data, err := operand.GetData()
	if err != nil {
		return nil, err
	}
	switch op := data.(type) {
	case string, float64, int, map[string]interface{}:
		return op, nil
	case *parser.Operation:
		inner, err := p.getMySQLQuery(op)
		if err != nil {
			return nil, err
		}
		return inner, nil
	case []parser.Operand:
		tmp := make([]interface{}, 0)
		for _, o := range op {
			inner, err := p.getMySQLOperands([]parser.Operand{o})
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

func insertNotOp(s interface{}) (string, error) {
	str, ok := s.(string)
	if !ok {
		return "", newUnsupportedOperandError(s)
	}
	parts := strings.Split(str, " ")
	if len(parts) >= 2 {
		return parts[0] + " NOT " + strings.Join(parts[1:], " "), nil
	}
	return "NOT " + str, nil
}

func (p *Parser) escapeColName(s interface{}) string {
	switch data := s.(type) {
	case string:
		if p.functionMatch.MatchString(data) || p.alreadyEscapedMatch.MatchString(data) {
			// We don't need to escape function names
			return data
		}
		return "`" + data + "`"
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	case int:
		return strconv.Itoa(data)
	default:
		return badString
	}
}

func escapeValue(s interface{}) string {
	switch data := s.(type) {
	case string:
		return "'" + data + "'"
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	case int:
		return strconv.Itoa(data)
	case []interface{}:
		ret := "("
		for i, v := range data {
			ret += escapeValue(v)
			if i != len(data)-1 {
				ret += ","
			}
		}
		ret += ")"
		return ret
	case map[string]interface{}:
		//nolint:errchkjson // This was unmarshaled from JSON, so it should be valid
		jsonData, _ := json.Marshal(data)
		return "'" + strings.ReplaceAll(string(jsonData), `"`, `\"`) + "'"
	default:
		return badString
	}
}

func escapeJSONValue(s interface{}) string {
	switch data := s.(type) {
	case string:
		return `"` + data + `"`
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	case int:
		return strconv.Itoa(data)
	case []interface{}:
		ret := "["
		for i, v := range data {
			ret += escapeJSONValue(v)
			if i != len(data)-1 {
				ret += ","
			}
		}
		ret += "]"
		return ret
	default:
		return badString
	}
}
