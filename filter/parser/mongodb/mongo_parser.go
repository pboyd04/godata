package mongodb

import (
	"strconv"

	"github.com/pboyd04/godata/filter/lexer"
	"github.com/pboyd04/godata/filter/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Parser struct {
}

func init() {
	// Register the parser
	parser.RegisterParser("mongodb", &Parser{})
}

func (p *Parser) GetDBQuery(common *parser.Parser) (interface{}, error) {
	op, err := common.GetOperation()
	if err != nil {
		return nil, err
	}
	return p.getMongoQuery(op)
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
	return p.getMongoQuery(op)
}

//nolint:funlen,cyclop
func (p *Parser) getMongoQuery(op *parser.Operation) (bson.D, error) {
	operands, err := p.getMongoOperands(op.Operands)
	if err != nil {
		return nil, err
	}
	if len(operands) > 1 {
		op, ok := operands[0].(string)
		if ok && op == "_id" {
			// convert the id to an oid doc...
			oid, err := primitive.ObjectIDFromHex(operands[1].(string))
			if err != nil {
				return nil, err
			}
			operands[1] = oid
		}
	}
	//nolint:exhaustive // This won't cover everything and will use the default case to catch errors
	switch op.Operator {
	case lexer.TokenTrue:
		return bson.D{}, nil
	case lexer.TokenFalse:
		return bson.D{{Key: "_id", Value: bson.D{{Key: "$type", Value: "string"}}}}, nil
	case lexer.Equals:
		return doSimpleOp("$eq", operands[0], operands[1])
	case lexer.NotEquals:
		return doSimpleOp("$ne", operands[0], operands[1])
	case lexer.GreaterThan:
		return doSimpleOp("$gt", operands[0], operands[1])
	case lexer.GreaterThanOrEqual:
		return doSimpleOp("$gte", operands[0], operands[1])
	case lexer.LessThan:
		return doSimpleOp("$lt", operands[0], operands[1])
	case lexer.LessThanOrEqual:
		return doSimpleOp("$lte", operands[0], operands[1])
	case lexer.And:
		return bson.D{{Key: "$and", Value: operands}}, nil
	case lexer.Or:
		return bson.D{{Key: "$or", Value: operands}}, nil
	case lexer.In:
		return doArrayOp("$in", operands[0], operands[1])
	case lexer.Contains:
		return doRegExOp("", "", operands[0], operands[1])
	case lexer.EndsWith:
		return doRegExOp("", "$", operands[0], operands[1])
	case lexer.Not:
		childOp, ok := operands[0].(bson.D)
		if !ok {
			return nil, newParserError("attempting to do a not on a non-document field")
		}
		return bson.D{{Key: childOp[0].Key, Value: bson.D{{Key: "$not", Value: childOp[0].Value}}}}, nil
	case lexer.Length:
		strOp0, ok := operands[0].(string)
		if !ok {
			return nil, newParserError("attempting to do a length on a non-string field")
		}
		arr := []interface{}{bson.D{{Key: "$strLenCP", Value: "$" + strOp0}}}
		return bson.D{{Key: "$expr", Value: arr}}, nil
	case lexer.StartsWith:
		return doRegExOp("^", "", operands[0], operands[1])
	case lexer.HasSubset:
		return doArrayOp("$all", operands[0], operands[1])
	default:
		return nil, newUnsupportedOperatorError(op.Operator)
	}
}

func doSimpleOp(key string, operand0, operand1 interface{}) (bson.D, error) {
	switch op0Data := operand0.(type) {
	case string:
		return bson.D{{Key: op0Data, Value: bson.D{{Key: key, Value: operand1}}}}, nil
	case bson.D:
		arr, ok := op0Data[0].Value.([]interface{})
		if !ok {
			return nil, newParserError("Attempting to do a simple operation on an expression where the expression is not in correct format")
		}
		op1Int, ok := operand1.(int)
		if !ok {
			return nil, newParserError("Attempting to do a simple operation on an expression with a non-int value")
		}
		arr = append(arr, bson.D{{Key: "$numberDecimal", Value: strconv.Itoa(op1Int)}})
		return bson.D{{Key: op0Data[0].Key, Value: bson.D{{Key: key, Value: arr}}}}, nil
	default:
		return nil, newParserError("Attempting to do a simple operation on an unknown field")
	}
}

func doArrayOp(key string, operand0, operand1 interface{}) (bson.D, error) {
	strOp0, ok := operand0.(string)
	if !ok {
		return nil, newParserError("Attempting to do an array op on a non-string field")
	}
	arrOp1, ok := operand1.([]interface{})
	if !ok {
		return nil, newParserError("Attempting to do an array op on a non-array value")
	}
	return bson.D{{Key: strOp0, Value: bson.D{{Key: key, Value: arrOp1}}}}, nil
}

func doRegExOp(prefix, postfix string, operand0, operand1 interface{}) (bson.D, error) {
	strOp0, ok := operand0.(string)
	if !ok {
		return nil, newParserError("Attempting to do a regex on a non-string field")
	}
	strOp1, ok := operand1.(string)
	if !ok {
		return nil, newParserError("Attempting to do a regex with a non-string value")
	}
	regEx := prefix + strOp1 + postfix
	return bson.D{{Key: strOp0, Value: bson.D{{Key: "$regex", Value: regEx}}}}, nil
}

func (p *Parser) getMongoOperands(operands []parser.Operand) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, operand := range operands {
		tmp, err := p.getMongoOperand(operand)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (p *Parser) getMongoOperand(operand parser.Operand) (interface{}, error) {
	data, err := operand.GetData()
	if err != nil {
		return nil, err
	}
	switch op := data.(type) {
	case string, float64, int, map[string]interface{}:
		return op, nil
	case *parser.Operation:
		inner, err := p.getMongoQuery(op)
		if err != nil {
			return nil, err
		}
		return inner, nil
	case []parser.Operand:
		tmp := make([]interface{}, 0)
		for _, o := range op {
			inner, err := p.getMongoOperands([]parser.Operand{o})
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
