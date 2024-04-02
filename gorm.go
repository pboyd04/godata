//go:build !no_gin
// +build !no_gin

package odata

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pboyd04/godata/orderby"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errQueryNotSupported = errors.New("query not incorrect format")

func GetGormSettingsFromGin(c *gin.Context, dbInput *gorm.DB) (*gorm.DB, error) {
	dbOut := dbInput
	queryOpts, ok := c.Value("odata").(*QueryOptions)
	if !ok {
		return dbInput, nil
	}
	if queryOpts.Top != 0 {
		dbOut = dbOut.Limit(int(queryOpts.Top))
	}
	if queryOpts.Skip != 0 {
		dbOut = dbOut.Offset(int(queryOpts.Skip))
	}
	if queryOpts.OrderBy != nil {
		for _, order := range queryOpts.OrderBy.OrderItem {
			dbOut = dbOut.Order(clause.OrderByColumn{Column: clause.Column{Name: order.Property}, Desc: (order.Direction == orderby.DESC)})
		}
	}
	if queryOpts.Filter != nil {
		myQuery, err := queryOpts.Filter.GetDBQuery("gorm")
		if err != nil {
			return nil, err
		}
		queryArgs, ok := myQuery.([]interface{})
		if ok {
			dbOut = dbOut.Where(queryArgs[0], queryArgs[1:]...)
		} else {
			return nil, errQueryNotSupported
		}
	}
	if queryOpts.Select != nil {
		dbOut = dbOut.Select(*queryOpts.Select)
	}
	return dbOut, nil
}
