// Package middleware provides an HTTP middleware for parsing OData query options from the URL and adding them to the request context.
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	odata "github.com/pboyd04/godata"
)

type contextKeyType string
type processingFn func(string, *odata.QueryOptions) error
type OdataMiddleware struct {
	handler                http.Handler
	preProcessingFunctions map[string]processingFn
}

const (
	ContextKey contextKeyType = "odata"
)

func NewOdataMiddleware(handlerToWrap http.Handler) *OdataMiddleware {
	ret := &OdataMiddleware{handler: handlerToWrap}
	ret.EnableFilterSupport()
	ret.EnableSelectSupport()
	ret.EnableOrderBySupport()
	ret.EnableTopSupport()
	ret.EnableSkipSupport()
	ret.EnableCountSupport()
	return ret
}

func (o *OdataMiddleware) EnableFilterSupport() {
	o.addPreProcessingFunction("$filter", processFilter)
}

func (o *OdataMiddleware) DisableFilterSupport() {
	delete(o.preProcessingFunctions, "$filter")
}

func (o *OdataMiddleware) EnableSelectSupport() {
	o.addPreProcessingFunction("$select", processSelect)
}

func (o *OdataMiddleware) DisableSelectSupport() {
	delete(o.preProcessingFunctions, "$select")
}

func (o *OdataMiddleware) EnableOrderBySupport() {
	o.addPreProcessingFunction("$orderby", processOrderBy)
}

func (o *OdataMiddleware) DisableOrderBySupport() {
	delete(o.preProcessingFunctions, "$orderby")
}

func (o *OdataMiddleware) EnableTopSupport() {
	o.addPreProcessingFunction("$top", processTop)
}

func (o *OdataMiddleware) DisableTopSupport() {
	delete(o.preProcessingFunctions, "$top")
}

func (o *OdataMiddleware) EnableSkipSupport() {
	o.addPreProcessingFunction("$skip", processSkip)
}

func (o *OdataMiddleware) DisableSkipSupport() {
	delete(o.preProcessingFunctions, "$skip")
}

func (o *OdataMiddleware) EnableCountSupport() {
	o.addPreProcessingFunction("$count", processCount)
}

func (o *OdataMiddleware) DisableCountSupport() {
	delete(o.preProcessingFunctions, "$count")
}

func (o *OdataMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	odata := odata.NewQueryOptions()
	for key, fn := range o.preProcessingFunctions {
		queryContent := r.URL.Query().Get(key)
		if queryContent != "" {
			err := fn(queryContent, odata)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	ctx := context.WithValue(r.Context(), ContextKey, odata)
	if o.handler != nil {
		o.handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetOdataFromContext(ctx context.Context) *odata.QueryOptions {
	ctxKey, ok := ctx.Value(ContextKey).(*odata.QueryOptions)
	if !ok {
		return nil
	}
	return ctxKey
}

func (o *OdataMiddleware) addPreProcessingFunction(name string, fn processingFn) {
	if o.preProcessingFunctions == nil {
		o.preProcessingFunctions = make(map[string]processingFn)
	}
	o.preProcessingFunctions[name] = fn
}

func processFilter(filter string, o *odata.QueryOptions) error {
	return o.AddFilter(filter)
}

func processSelect(selectString string, o *odata.QueryOptions) error {
	o.AddSelect(strings.Split(selectString, ","))
	return nil
}

func processOrderBy(orderBy string, o *odata.QueryOptions) error {
	return o.AddOrderBy(orderBy)
}

func processTop(top string, o *odata.QueryOptions) error {
	topInt, err := strconv.ParseInt(top, 10, 64)
	if err != nil {
		return err
	}
	o.AddTop(topInt)
	return nil
}

func processSkip(skip string, o *odata.QueryOptions) error {
	skipInt, err := strconv.ParseInt(skip, 10, 64)
	if err != nil {
		return err
	}
	o.AddSkip(skipInt)
	return nil
}

func processCount(count string, o *odata.QueryOptions) error {
	countBool, err := strconv.ParseBool(count)
	if err != nil {
		return err
	}
	o.AddCount(countBool)
	return nil
}
