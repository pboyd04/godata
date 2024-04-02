package orderby

import "strings"

type OrderDirection int

const (
	ASC  OrderDirection = 1
	DESC OrderDirection = 2
)

type OrderBy struct {
	OrderItem []OrderItem
}

type OrderItem struct {
	Property  string
	Direction OrderDirection
}

func NewOrderBy(orderByString string) (*OrderBy, error) {
	orderBy := &OrderBy{}
	orderBy.OrderItem = make([]OrderItem, 0)
	return orderBy, orderBy.parseOrderBy(orderByString)
}

func (o *OrderBy) parseOrderBy(orderByString string) error {
	if len(orderByString) == 0 {
		return nil
	}
	orderByItems := strings.Split(orderByString, ",")
	for _, orderByItem := range orderByItems {
		orderItem, err := o.parseOrderItem(orderByItem)
		if err != nil {
			return err
		}
		o.OrderItem = append(o.OrderItem, orderItem)
	}
	return nil
}

func (o *OrderBy) parseOrderItem(orderByItem string) (OrderItem, error) {
	parts := strings.Split(orderByItem, " ")
	if len(parts) != 2 {
		return OrderItem{Property: parts[0], Direction: ASC}, nil
	}
	direction := strings.ToUpper(parts[1])
	switch direction {
	case "ASC":
		return OrderItem{Property: parts[0], Direction: ASC}, nil
	case "DESC":
		return OrderItem{Property: parts[0], Direction: DESC}, nil
	default:
		return OrderItem{}, &InvalidOrderDirectionError{Direction: direction}
	}
}

type InvalidOrderDirectionError struct {
	Direction string
}

func (e *InvalidOrderDirectionError) Error() string {
	return "Invalid order direction: " + e.Direction
}
