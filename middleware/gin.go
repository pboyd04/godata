//go:build !no_gin
// +build !no_gin

package middleware

import (
	"github.com/gin-gonic/gin"
	odata "github.com/pboyd04/godata"
)

func (o *OdataMiddleware) GinMiddleware(c *gin.Context) {
	queryOptions := odata.QueryOptions{}
	for k, v := range c.Request.URL.Query() {
		if fn, ok := o.preProcessingFunctions[k]; ok {
			if err := fn(v[0], &queryOptions); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	}
	c.Set(string(ContextKey), &queryOptions)
	c.Next()
	// Post processing
}
