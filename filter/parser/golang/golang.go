package golang

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/pboyd04/godata/filter/lexer"
	"github.com/pboyd04/godata/filter/parser"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Register the parser
	parser.RegisterParser("golang", &Parser{})
}

type Parser struct {
}

type internalValueState struct {
	value                interface{}
	currentComputedValue map[string]interface{}
	constant             interface{}
	computedConstant     interface{}
	isNilConstant        bool
}

type Evaluator struct {
	op *parser.Operation
}

func (p *Parser) GetDBQuery(common *parser.Parser) (interface{}, error) {
	op, err := common.GetOperation()
	if err != nil {
		return nil, err
	}
	ret := new(Evaluator)
	ret.op = op
	return ret, nil
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
	ret := new(Evaluator)
	ret.op = op
	return ret, nil
}

func newInternalValueState(data []interface{}) []internalValueState {
	ret := make([]internalValueState, len(data))
	for i, val := range data {
		myInternalState := internalValueState{value: val}
		myType := reflect.TypeOf(val)
		myInternalState.currentComputedValue = make(map[string]interface{})
		if myType.Kind() == reflect.Struct {
			myFields := reflect.VisibleFields(myType)
			for _, field := range myFields {
				name := field.Name
				if field.Tag.Get("json") != "" {
					name = field.Tag.Get("json")
				}
				myInternalState.currentComputedValue[name] = reflect.ValueOf(val).FieldByIndex(field.Index).Interface()
			}
		} else if myType.Kind() == reflect.Map {
			myMap, ok := val.(map[string]interface{})
			if !ok {
				log.Printf("Unknown map type %T\n", val)
				continue
			}
			for k, v := range myMap {
				myInternalState.currentComputedValue[k] = v
			}
		}
		ret[i] = myInternalState
	}
	return ret
}

type dummyTesting struct{}

func (d dummyTesting) Errorf(_ string, _ ...interface{}) {}

type comparisonFn func(a, b interface{}) bool
type opFunction func(val interface{}, state *internalValueState, additionalOperands ...*internalValueState) interface{}

//nolint:gochecknoglobals // This is a map of functions that are used to compare values it's easier to keep it this way than to have a bunch of if statements
var comparisonMap = map[lexer.TokenKey]comparisonFn{
	lexer.Equals: func(a, b interface{}) bool {
		// I don't really get it, but nil is not always nil (i.e. if there is type associated with it, this is basically the only fool proof way that doesn't panic)
		if (*[2]uintptr)(unsafe.Pointer(&a))[1] == 0 && (*[2]uintptr)(unsafe.Pointer(&b))[1] == 0 {
			return true
		}
		switch aVal := a.(type) {
		case float64:
			_, ok := b.(float64)
			if !ok {
				bVal, ok := b.(int)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return int(aVal) == bVal
			}
			return a == b
		default:
			log.Printf("%#v == %#v\n", a, b)
			return a == b
		}
	},
	lexer.NotEquals: func(a, b interface{}) bool {
		return a != b
	},
	lexer.GreaterThan: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case int:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return float64(aVal) > bVal
			}
			return aVal > bVal
		case float64:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return aVal > bVal
			}
			return aVal > float64(bVal)
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
				return false
			}
			return strings.Compare(aVal, bVal) > 0
		default:
			log.Printf("Unknown compare for primary type %T\n", a)
			return false
		}
	},
	lexer.GreaterThanOrEqual: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case int:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return float64(aVal) >= bVal
			}
			return aVal >= bVal
		case float64:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return aVal >= bVal
			}
			return aVal >= float64(bVal)
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
				return false
			}
			return strings.Compare(aVal, bVal) >= 0
		default:
			log.Printf("Unknown compare for primary type %T\n", a)
			return false
		}
	},
	lexer.LessThan: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case int:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return float64(aVal) < bVal
			}
			return aVal < bVal
		case float64:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return aVal < bVal
			}
			return aVal < float64(bVal)
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
				return false
			}
			return strings.Compare(aVal, bVal) < 0
		default:
			log.Printf("Unknown compare for primary type %T\n", a)
			return false
		}
	},
	lexer.LessThanOrEqual: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case int:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return float64(aVal) <= bVal
			}
			return aVal <= bVal
		case float64:
			bVal, ok := b.(int)
			if !ok {
				bVal, ok := b.(float64)
				if !ok {
					log.Printf("Unknown numeric compare for second type %T\n", b)
					return false
				}
				return aVal <= bVal
			}
			return aVal <= float64(bVal)
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
				return false
			}
			return strings.Compare(aVal, bVal) <= 0
		default:
			log.Printf("Unknown compare for primary type %T\n", a)
			return false
		}
	},
	lexer.Contains: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
			}
			return strings.Contains(aVal, bVal)
		default:
			log.Printf("Unknown contains compare for primary type %T\n", a)
			return false
		}
	},
	lexer.EndsWith: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
			}
			return strings.HasSuffix(aVal, bVal)
		default:
			log.Printf("Unknown endswith compare for primary type %T\n", a)
			return false
		}
	},
	lexer.StartsWith: func(a, b interface{}) bool {
		switch aVal := a.(type) {
		case string:
			bVal, ok := b.(string)
			if !ok {
				log.Printf("Unknown string compare for second type %T\n", b)
			}
			return strings.HasPrefix(aVal, bVal)
		default:
			log.Printf("Unknown startswith compare for primary type %T\n", a)
			return false
		}
	},
	lexer.HasSubset: func(a, b interface{}) bool {
		switch bData := b.(type) {
		case []parser.Operand:
			bVal := make([]interface{}, 0)
			for _, val := range bData {
				valData, err := val.GetData()
				if err != nil {
					log.Printf("Error getting data from operand: %v\n", err)
					return false
				}
				bVal = append(bVal, valData)
			}
			return assert.Subset(dummyTesting{}, a, bVal)
		default:
			return assert.Subset(dummyTesting{}, a, b)
		}
	},
	lexer.HasSubsequence: func(a, b interface{}) bool {
		var aVal []interface{}
		var bVal []interface{}
		switch val := a.(type) {
		case []string:
			aVal = make([]interface{}, len(val))
			for i, v := range val {
				aVal[i] = v
			}
		case []int:
			aVal = make([]interface{}, len(val))
			for i, v := range val {
				aVal[i] = v
			}
		case []float64:
			aVal = make([]interface{}, len(val))
			for i, v := range val {
				aVal[i] = v
			}
		default:
			log.Printf("Unknown has subsequence compare for primary type %T\n", val)
			return false
		}
		switch val := b.(type) {
		case []interface{}:
			bVal = val
		default:
			log.Printf("Unknown has subsequence compare for secondary type %T\n", val)
			return false
		}
		return hasSubsequence(aVal, bVal)
	},
	lexer.MatchesPattern: func(a, b interface{}) bool {
		switch v := a.(type) {
		case string:
			reg, err := regexp.Compile(b.(string))
			if err != nil {
				log.Printf("Error compiling regex: %v\n", err)
				return false
			}
			return reg.MatchString(v)
		default:
			log.Printf("Unknown matches pattern for primary type %T %#v\n", v, b)
			return false
		}
	},
}

func hasSubsequence(a []interface{}, b []interface{}) bool {
	if len(b) == 0 {
		return true
	}
	if len(b) > len(a) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] == b[0] {
			if len(b) == 1 {
				return true
			}
			if i+1 < len(a) {
				if hasSubsequence(a[i+1:], b[1:]) {
					return true
				}
			}
		}
	}
	return false
}

//nolint:gochecknoglobals // This is a map of functions that are used to operate on values it's easier to keep it this way than to have a bunch of if statements
var opMap = map[lexer.TokenKey]opFunction{
	lexer.Length: func(val interface{}, _ *internalValueState, _ ...*internalValueState) interface{} {
		switch v := val.(type) {
		case string:
			return len(v)
		default:
			log.Printf("Unknown length compare for primary type %T %v\n", val, val)
			return 0
		}
	},
	lexer.Add: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Add requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			if additionalOperands[0].IsFloat() {
				return decimal.NewFromInt(int64(v)).Add(additionalOperands[0].Float())
			}
			return v + additionalOperands[0].Int()
		case float64:
			ret, _ := decimal.NewFromFloat(v).Add(additionalOperands[0].Float()).Float64()
			return ret
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Subtract: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Subtract requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			if additionalOperands[0].IsFloat() {
				return decimal.NewFromInt(int64(v)).Sub(additionalOperands[0].Float())
			}
			return v - additionalOperands[0].Int()
		case float64:
			ret, _ := decimal.NewFromFloat(v).Sub(additionalOperands[0].Float()).Float64()
			return ret
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Multiply: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Multiply requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			if additionalOperands[0].IsFloat() {
				return decimal.NewFromInt(int64(v)).Mul(additionalOperands[0].Float())
			}
			return v * additionalOperands[0].Int()
		case float64:
			ret, _ := decimal.NewFromFloat(v).Mul(additionalOperands[0].Float()).Float64()
			return ret
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Divide: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Divide requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			if additionalOperands[0].IsFloat() {
				ret, _ := decimal.NewFromInt(int64(v)).Div(additionalOperands[0].Float()).Float64()
				return ret
			}
			return v / additionalOperands[0].Int()
		case float64:
			ret, _ := decimal.NewFromFloat(v).Div(additionalOperands[0].Float()).Float64()
			return ret
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.DivideFloat: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Divide requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			ret, _ := decimal.NewFromInt(int64(v)).Div(additionalOperands[0].Float()).Float64()
			return ret
		case float64:
			ret, _ := decimal.NewFromFloat(v).Div(additionalOperands[0].Float()).Float64()
			return ret
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Modulo: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Divide requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case int:
			return v % additionalOperands[0].Int()
		default:
			log.Printf("Unknown add compare for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Concat: func(val interface{}, state *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Concat requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case string:
			ret := v + additionalOperands[0].String(state)
			return ret
		default:
			log.Printf("Unknown concat for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.IndexOf: func(val interface{}, state *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("IndexOf requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case string:
			ret := strings.Index(v, additionalOperands[0].String(state))
			return ret
		default:
			log.Printf("Unknown indexof for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Substring: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		if len(additionalOperands) < 1 {
			log.Printf("Substring requires at least one operand\n")
			return 0
		}
		switch v := val.(type) {
		case string:
			if len(additionalOperands) == 1 {
				ret := v[additionalOperands[0].Int():]
				return ret
			}
			start := additionalOperands[0].Int()
			end := additionalOperands[1].Int()
			end = start + end
			length := len(v)
			if end > length {
				end = length - 1
			}
			ret := v[start:end]
			return ret
		default:
			log.Printf("Unknown substring for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.ToLower: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case string:
			return strings.ToLower(v)
		default:
			log.Printf("Unknown tolower for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.ToUpper: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case string:
			return strings.ToUpper(v)
		default:
			log.Printf("Unknown toupper for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Trim: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case string:
			return strings.TrimSpace(v)
		default:
			log.Printf("Unknown trim for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Day: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return v.Day()
		default:
			log.Printf("Unknown day for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.FractionalSeconds: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			ret := float64(v.Nanosecond()) / float64(1000000)
			return ret
		default:
			log.Printf("Unknown fractional seconds for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Hour: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return v.Hour()
		default:
			log.Printf("Unknown hour for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Minute: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return v.Minute()
		default:
			log.Printf("Unknown minute for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Month: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return int(v.Month())
		default:
			log.Printf("Unknown month for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Second: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return v.Second()
		default:
			log.Printf("Unknown second for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Year: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case time.Time:
			return v.Year()
		default:
			log.Printf("Unknown year for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Ceiling: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case float64:
			ret, _ := decimal.NewFromFloat(v).Ceil().Float64()
			return ret
		default:
			log.Printf("Unknown ceiling for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Floor: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case float64:
			ret, _ := decimal.NewFromFloat(v).Floor().Float64()
			return ret
		default:
			log.Printf("Unknown floor for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
	lexer.Round: func(val interface{}, _ *internalValueState, additionalOperands ...*internalValueState) interface{} {
		switch v := val.(type) {
		case float64:
			ret, _ := decimal.NewFromFloat(v).Round(0).Float64()
			return ret
		default:
			log.Printf("Unknown round for primary type %T %#v\n", v, additionalOperands)
			return 0
		}
	},
}

func (e *Evaluator) FilterSlice(data []interface{}) ([]interface{}, error) {
	if len(data) == 0 {
		return data, nil
	}
	state := newInternalValueState(data)
	return e.filterSlice(state, e.op)
}

func (e *Evaluator) filterSlice(data []internalValueState, op *parser.Operation) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, d := range data {
		ok, err := d.passesOp(op)
		if err != nil {
			return nil, err
		}
		if ok {
			ret = append(ret, d.value)
		}
	}
	return ret, nil
}

//nolint:cyclop // This function is complex because it has to handle all the different types of operations
func (d *internalValueState) passesOp(op *parser.Operation) (bool, error) {
	operands, err := d.getStatesFromOperands(op.Operands)
	if err != nil {
		return false, err
	}
	switch op.Operator {
	case lexer.TokenTrue:
		return true, nil
	case lexer.TokenFalse:
		return false, nil
	case lexer.Equals, lexer.NotEquals, lexer.GreaterThan, lexer.GreaterThanOrEqual, lexer.LessThan, lexer.LessThanOrEqual, lexer.Contains, lexer.EndsWith, lexer.StartsWith, lexer.HasSubset, lexer.HasSubsequence, lexer.MatchesPattern:
		return d.simpleCompare(operands[0], operands[1], comparisonMap[lexer.TokenKey(op.Operator)])
	case lexer.And:
		if operands[0] == nil || operands[1] == nil {
			return false, nil
		}
		return true, nil
	case lexer.Or:
		if operands[0] == nil && operands[1] == nil {
			return false, nil
		}
		return true, nil
	case lexer.In:
		return d.in(operands[0], operands[1])
	case lexer.Not:
		return operands[0] == nil, nil
	case lexer.Length, lexer.Add, lexer.Subtract, lexer.Multiply, lexer.Divide, lexer.DivideFloat, lexer.Modulo, lexer.Concat, lexer.IndexOf, lexer.Substring, lexer.ToLower, lexer.ToUpper, lexer.Trim, lexer.Day, lexer.FractionalSeconds, lexer.Hour, lexer.Minute, lexer.Month, lexer.Second, lexer.Year, lexer.Ceiling, lexer.Floor, lexer.Round:
		return d.computeOperation(operands, lexer.TokenKey(op.Operator))
	case parser.NoOp:
		return true, nil
	default:
		return false, &UnsupportedOperatorError{operator: lexer.TokenKey(op.Operator)}
	}
}

func (d *internalValueState) IsFloat() bool {
	_, ok := d.computedConstant.(float64)
	if !ok && d.constant != nil {
		_, ok = d.constant.(float64)
	}
	return ok
}

func (d *internalValueState) Float() decimal.Decimal {
	value := d.computedConstant
	if value == nil {
		value = d.constant
	}
	switch val := value.(type) {
	case float64:
		return decimal.NewFromFloat(val)
	case int:
		return decimal.NewFromInt(int64(val))
	default:
		log.Printf("Unknown float type %T\n", val)
		return decimal.NewFromInt(0)
	}
}

func (d *internalValueState) Int() int {
	value := d.computedConstant
	if value == nil {
		value = d.constant
	}
	switch val := value.(type) {
	case float64:
		return int(val)
	case int:
		return val
	default:
		log.Printf("Unknown int type %T\n", val)
		return 0
	}
}

func (d *internalValueState) String(state *internalValueState) string {
	value := d.computedConstant
	if value == nil {
		value = d.constant
	}
	strValue, ok := value.(string)
	if ok {
		fieldValue, ok := state.currentComputedValue[strValue]
		if ok {
			value = fieldValue
		}
	}
	switch val := value.(type) {
	case string:
		return val
	default:
		log.Printf("Unknown string type %T\n", val)
		return ""
	}
}

func (d *internalValueState) simpleCompare(op1 *internalValueState, op2 *internalValueState, compareFn comparisonFn) (bool, error) {
	var value interface{}
	var value2 interface{}
	if op1.computedConstant != nil {
		value = op1.computedConstant
	} else {
		strVal, ok := op1.constant.(string)
		if !ok {
			return false, &UnsupportedDataTypeError{}
		}
		fieldValue, ok := d.currentComputedValue[strVal]
		if !ok {
			return false, &UnknownFieldError{field: strVal}
		}
		value = fieldValue
	}
	switch {
	case op2.isNilConstant:
		value2 = nil
	case op2.computedConstant != nil:
		value2 = op2.computedConstant
	default:
		value2 = op2.constant
		strValue2, ok := op2.constant.(string)
		if ok {
			fieldValue, ok := d.currentComputedValue[strValue2]
			if ok {
				value2 = fieldValue
			}
		}
	}
	return compareFn(value, value2), nil
}

func (d *internalValueState) in(op1 *internalValueState, op2 *internalValueState) (bool, error) {
	strVal, ok := op1.constant.(string)
	if !ok {
		return false, &UnsupportedDataTypeError{}
	}
	fieldValue, ok := d.currentComputedValue[strVal]
	if !ok {
		return false, &UnknownFieldError{field: strVal}
	}
	return assert.Contains(dummyTesting{}, op2.constant, fieldValue), nil
}

func (d *internalValueState) computeOperation(operands []*internalValueState, op lexer.TokenKey) (bool, error) {
	if len(operands) < 1 {
		return false, newParserError("computed operations require at least one operand")
	}
	var value interface{}
	if operands[0].computedConstant != nil {
		value = operands[0].computedConstant
	} else {
		strVal, ok := operands[0].constant.(string)
		if !ok {
			return false, &UnsupportedDataTypeError{}
		}
		fieldValue, ok := d.currentComputedValue[strVal]
		if !ok {
			return false, &UnknownFieldError{field: strVal}
		}
		value = fieldValue
	}
	opMapFn, ok := opMap[op]
	if !ok {
		return false, &UnsupportedOperatorError{operator: op}
	}
	d.computedConstant = opMapFn(value, d, operands[1:]...)
	return true, nil
}

func (d *internalValueState) getStatesFromOperands(operands []parser.Operand) ([]*internalValueState, error) {
	ret := make([]*internalValueState, len(operands))
	for i, operand := range operands {
		tmp, err := d.getStateFromOperand(operand)
		if err != nil {
			return nil, err
		}
		ret[i] = tmp
	}
	return ret, nil
}

//nolint:cyclop // This function isn't that complex
func (d *internalValueState) getStateFromOperand(operand parser.Operand) (*internalValueState, error) {
	data, err := operand.GetData()
	if err != nil {
		return nil, err
	}
	switch op := data.(type) {
	case string, float64, int:
		return &internalValueState{constant: op}, nil
	case *parser.Operation:
		passes, err := d.passesOp(op)
		if err != nil {
			return nil, err
		}
		if !passes {
			return nil, nil
		}
		return d, nil
	case []parser.Operand:
		ret := new(internalValueState)
		ret.constant = make([]interface{}, 0)
		for _, o := range op {
			inner, err := d.getStateFromOperand(o)
			if err != nil {
				return nil, err
			}
			//nolint:forcetypeassert // The variable was assigned just above
			ret.constant = append(ret.constant.([]interface{}), inner.constant)
		}
		return ret, nil
	case nil:
		return &internalValueState{isNilConstant: true}, nil
	default:
		return nil, newUnsupportedOperandError(op)
	}
}
