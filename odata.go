package odata

import (
	"github.com/pboyd04/godata/filter"
	"github.com/pboyd04/godata/orderby"
)

type QueryOptions struct {
	Filter  *filter.Filter
	Select  *[]string
	OrderBy *orderby.OrderBy
	Top     int64
	Skip    int64
	Count   bool
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{Top: -1, Skip: -1, Count: false}
}

func (q *QueryOptions) AddFilter(filterString string) error {
	f, err := filter.NewFilter(filterString)
	if err != nil {
		return err
	}
	q.Filter = f
	return nil
}

func (q *QueryOptions) AddSelect(selects []string) {
	if q.Select == nil {
		q.Select = &selects
		return
	}
	*q.Select = append(*q.Select, selects...)
}

func (q *QueryOptions) AddOrderBy(orderByString string) error {
	o, err := orderby.NewOrderBy(orderByString)
	if err != nil {
		return err
	}
	q.OrderBy = o
	return nil
}

func (q *QueryOptions) AddTop(top int64) {
	q.Top = top
}

func (q *QueryOptions) AddSkip(skip int64) {
	q.Skip = skip
}

func (q *QueryOptions) AddCount(count bool) {
	q.Count = count
}
