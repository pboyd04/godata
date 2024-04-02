package parser

import (
	"errors"

	"github.com/pboyd04/godata/filter/lexer"
)

type tokenGroup struct {
	children []Operand
}

func (t *tokenGroup) GetData() (interface{}, error) {
	return t, nil
}

func newTokenGroup(tokens []*lexer.Token) *tokenGroup {
	ret := &tokenGroup{}
	ret.children = make([]Operand, len(tokens))
	// Not sure why, but append doesn't work in this case
	for i, token := range tokens {
		ret.children[i] = token
	}
	return ret
}

func newTokenGroupOperands(ops []Operand) *tokenGroup {
	ret := &tokenGroup{}
	ret.children = make([]Operand, len(ops))
	copy(ret.children, ops)
	return ret
}

func (t *tokenGroup) beforeAndAfterOperandProcess(operand Operand, compareFn tokenComparer, i int) error {
	switch token := operand.(type) {
	case *lexer.Token:
		if compareFn(token) {
			// Found an operator
			err := t.doBeforeAndAfterTokenReplace(token, i)
			if err != nil {
				return err
			}
		}
	case *tokenGroup:
		err := token.beforeAndAfterOpProcess(compareFn)
		if err != nil {
			return err
		}
	case *Operation:
		err := token.beforeAndAfterOpProcess(compareFn)
		if err != nil {
			return err
		}
	case *ObjectOperand, *SliceOperand:
		// Don't do anything this is already processed
		break
	default:
		return newParserError("unknown type: %T", token)
	}
	return nil
}

func (t *tokenGroup) beforeAndAfterOpProcess(compareFn tokenComparer) error {
	for i := 0; i < len(t.children); i++ {
		err := t.beforeAndAfterOperandProcess(t.children[i], compareFn, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tokenGroup) createSlices() {
	for i := 0; i < len(t.children); i++ {
		switch token := t.children[i].(type) {
		case *lexer.Token:
			if token.Type == lexer.OpenSquareBracket {
				for j := i + 1; j < len(t.children); j++ {
					token, ok := t.children[j].(*lexer.Token)
					if ok && token.Type == lexer.CloseSquareBracket {
						// Found the end of the group
						op := newSliceOperand(t.children[i+1 : j])
						tmp := t.children[j+1:]
						t.children = append(t.children[:i], op)
						t.children = append(t.children, tmp...)
						break
					}
				}
			}
		default:
			break
		}
	}
}

func (t *tokenGroup) doBeforeAndAfterTokenReplace(token *lexer.Token, i int) error {
	op := &Operation{Operator: Operator(token.Type)}
	// Previous item and next item must both be something... we don't care what...
	if i == 0 || i == len(t.children)-1 {
		return newParserError("expected item before and after operator")
	}
	prev := t.children[i-1]
	next := t.children[i+1]
	op.Operands = []Operand{prev, next}
	tmp := t.children[i+2:]
	t.children = append(t.children[:i-1], op)
	t.children = append(t.children, tmp...)
	return nil
}

func (t *tokenGroup) findEndParentheses(start int) (int, bool) {
	depth := 1
	foundNested := false
	for i := start + 1; i < len(t.children); i++ {
		token, ok := t.children[i].(*lexer.Token)
		if ok && token.Type == lexer.OpenParens {
			depth++
			foundNested = true
		} else if ok && token.Type == lexer.CloseParens {
			depth--
			if depth == 0 {
				return i, foundNested
			}
		}
	}
	return -1, false
}

func (t *tokenGroup) flatten(parent *Operation) (*Operation, error) {
	if len(t.children) != 1 {
		return nil, ErrMoreThanOneChild
	}
	var op *Operation
	switch operand := t.children[0].(type) {
	case *Operation:
		op = operand
	case *tokenGroup:
		myOp, err := operand.toOperation(parent)
		if err != nil {
			return nil, err
		}
		op = myOp
	case *lexer.Token:
		return &Operation{Operator: Operator(operand.Type)}, nil
	default:
		return nil, newParserError("unknown type: %T", t.children[0])
	}
	err := op.flatten()
	if err != nil {
		return nil, err
	}
	return op, nil
}

func (t *tokenGroup) handleFlattenError(o *Operation, err error) error {
	if errors.Is(err, ErrMoreThanOneChild) {
		if o == nil {
			return newParserError("more than one non-operation child at root")
		}
		o.Operands = append(o.Operands, t.children...)
		// Remove all token groups from the parent
		for i := 0; i < len(o.Operands); i++ {
			_, ok := o.Operands[i].(*tokenGroup)
			if ok {
				o.Operands = append(o.Operands[:i], o.Operands[i+1:]...)
			}
		}
		return nil
	}
	return err
}

func (t *tokenGroup) handleInToken(i int) error {
	op := &Operation{Operator: lexer.In}
	// We need two operands, the first is before the in, the second is after the in
	prev := t.children[i-1]
	// There are two options, option 1 the next item is a token group, option 2 the next item is in square brackets...
	switch next := t.children[i+1].(type) {
	case *tokenGroup:
		// Found a token group
		err := next.in()
		if err != nil {
			return err
		}
		op.Operands = []Operand{prev, newSliceOperand(next.children)}
		tmp := t.children[i+2:]
		t.children = append(t.children[:i-1], op)
		t.children = append(t.children, tmp...)
	case *lexer.Token:
		if next.Type != lexer.OpenSquareBracket {
			// Can't find arguments for in
			return newParserError("expected open square bracket after in, found %s", next.Text)
		}
		for j := i + 2; j < len(t.children); j++ {
			token, ok := t.children[j].(*lexer.Token)
			if ok && token.Type == lexer.CloseSquareBracket {
				// Found the end of the group
				op.Operands = []Operand{prev, newSliceOperand(t.children[i+2 : j])}
				tmp := t.children[j+1:]
				t.children = append(t.children[:i-1], op)
				t.children = append(t.children, tmp...)
				break
			}
		}
	default:
		return newParserError("unknown type (inner loop): %T", t.children[i+1])
	}
	return nil
}

func (t *tokenGroup) handleMethodToken(token *lexer.Token, i int) error {
	op := &Operation{Operator: Operator(token.Type)}
	// The next item must be a token group otherwise it's an error
	next := t.children[i+1]
	group, ok := next.(*tokenGroup)
	if !ok {
		return newParserError("expected token group after method call, found %#v", next)
	}
	op.Operands = []Operand{group}
	tmp := t.children[i+2:]
	t.children = append(t.children[:i], op)
	t.children = append(t.children, tmp...)
	err := group.methodCalls()
	if err != nil {
		return err
	}
	group.createSlices()
	group.removeCommas()
	return nil
}

func (t *tokenGroup) handleObjects() error {
	var op *ObjectOperand
	initialIndex := -1
	for i := 0; i < len(t.children); i++ {
		token, ok := t.children[i].(*lexer.Token)
		if !ok {
			tg, ok := t.children[i].(*tokenGroup)
			if !ok {
				return newParserError("unknown type: %T", t.children[i])
			}
			err := tg.handleObjects()
			if err != nil {
				return err
			}
			continue
		}
		if token.Type == lexer.OpenCurlyBrace {
			// Found an object
			op = &ObjectOperand{Properties: ""}
			initialIndex = i
		}
		if op != nil {
			op.Properties += token.Text
		}
		if token.Type == lexer.CloseCurlyBrace {
			// Found the end of the object
			t.children[initialIndex] = op
			// Remove the tokens that were used to create the object
			t.children = append(t.children[:initialIndex+1], t.children[i+1:]...)
			op = nil
		}
	}
	return nil
}

func (t *tokenGroup) handleUnaryToken(token *lexer.Token, i int) error {
	op := &Operation{Operator: Operator(token.Type)}
	// The next item can be whatever...
	next := t.children[i+1]
	op.Operands = []Operand{next}
	tmp := t.children[i+2:]
	t.children = append(t.children[:i], op)
	switch nextVal := next.(type) {
	case *tokenGroup:
		err := nextVal.unary()
		if err != nil {
			return err
		}
		nextVal.children = append(nextVal.children, tmp...)
	case *Operation:
		err := nextVal.unary()
		if err != nil {
			return err
		}
		nextVal.Operands = append(nextVal.Operands, tmp...)
	case *lexer.Token:
		group := newTokenGroupOperands([]Operand{next})
		group.children = append(group.children, tmp...)
		err := group.unary()
		if err != nil {
			return err
		}
		// Replace the token with the group
		op.Operands = []Operand{group}
	}
	return nil
}

func (t *tokenGroup) in() error {
	// Handle the in operator
	for i := 0; i < len(t.children); i++ {
		switch token := t.children[i].(type) {
		case *lexer.Token:
			if token.Type == lexer.In {
				// Found an in operator
				err := t.handleInToken(i)
				if err != nil {
					return err
				}
				break
			}
		case *tokenGroup:
			err := token.in()
			if err != nil {
				return err
			}
		default:
			// Don't do anything this is already processed
		}
	}
	return nil
}

func (t *tokenGroup) methodCalls() error {
	for i := 0; i < len(t.children); i++ {
		switch token := t.children[i].(type) {
		case *lexer.Token:
			if token.HasParameters() {
				// Found a method call
				err := t.handleMethodToken(token, i)
				if err != nil {
					return err
				}
			}
		case *tokenGroup:
			err := token.methodCalls()
			if err != nil {
				return err
			}
		case *Operation:
			err := token.methodCalls()
			if err != nil {
				return err
			}
		default:
			// Don't do anything this is already processed
		}
	}
	return nil
}

func (t *tokenGroup) parentheses() {
	// Special case, this is always first so everything in the array is a token
	for i := 0; i < len(t.children); i++ {
		//nolint:forcetypeassert // We know this is a token
		token := t.children[i].(*lexer.Token)
		if token.Type == lexer.OpenParens {
			end, foundNested := t.findEndParentheses(i)
			if end == -1 {
				// No matching close parentheses
				return
			}
			newGroup := newTokenGroupOperands(t.children[i+1 : end])
			if foundNested {
				// Handle the nested parentheses
				newGroup.parentheses()
			}
			tmp := t.children[end+1:]
			t.children = append(t.children[:i], newGroup)
			t.children = append(t.children, tmp...)
		}
	}
}

func (t *tokenGroup) processTokenGroupOps() error {
	err := t.in()
	if err != nil {
		return err
	}
	err = t.unary()
	if err != nil {
		return err
	}
	err = t.methodCalls()
	if err != nil {
		return err
	}
	err = t.beforeAndAfterOpProcess(multiplicativeCompare)
	if err != nil {
		return err
	}
	err = t.beforeAndAfterOpProcess(additiveCompare)
	if err != nil {
		return err
	}
	err = t.beforeAndAfterOpProcess(relationalCompare)
	if err != nil {
		return err
	}
	err = t.beforeAndAfterOpProcess(equalityCompare)
	if err != nil {
		return err
	}
	err = t.beforeAndAfterOpProcess(conjunctionCompare)
	if err != nil {
		return err
	}
	return nil
}

func (t *tokenGroup) removeCommas() {
	// Remove commas from the token group
	for i := 0; i < len(t.children); i++ {
		token, ok := t.children[i].(*lexer.Token)
		if ok && token.Type == lexer.Comma {
			// Remove the comma
			t.children = append(t.children[:i], t.children[i+1:]...)
			i--
		}
	}
}

func (t *tokenGroup) toOperation(parent *Operation) (*Operation, error) {
	if parent != nil && parent.Operator.hasParameters() {
		// This is a method call
		parent.Operands = append(parent.Operands, t.children...)
		// Remove all token groups from the parent
		for i := 0; i < len(parent.Operands); i++ {
			_, ok := parent.Operands[i].(*tokenGroup)
			if ok {
				parent.Operands = append(parent.Operands[:i], parent.Operands[i+1:]...)
			}
		}
		err := parent.flatten()
		return nil, err
	}
	myOp, err := t.flatten(parent)
	if err != nil {
		return nil, t.handleFlattenError(parent, err)
	}
	return myOp, nil
}

func (t *tokenGroup) unary() error {
	for i := 0; i < len(t.children); i++ {
		switch token := t.children[i].(type) {
		case *lexer.Token:
			if token.IsUnary() {
				// Found a unary operator
				err := t.handleUnaryToken(token, i)
				if err != nil {
					return err
				}
			}
		case *tokenGroup:
			err := token.unary()
			if err != nil {
				return err
			}
		case *Operation:
			err := token.unary()
			if err != nil {
				return err
			}
		default:
			// Don't do anything this is a type that is already processed
		}
	}
	return nil
}
