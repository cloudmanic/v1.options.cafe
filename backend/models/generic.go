//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//
// A generic way to query any model we want.
//
func (t *DB) Query(model interface{}, params QueryParam) error {

	var query *gorm.DB

	// Validate order column
	if (len(params.Order) > 0) && (len(params.AllowedOrderCols) > 0) {
		var found = false

		for _, row := range params.AllowedOrderCols {
			if params.Order == row {
				found = true
			}
		}

		if !found {
			return errors.New("Invalid order parameter. - " + params.Order)
		}
	}

	// Do some quick filtering - Think injections
	var sortText = strings.ToUpper(params.Sort)
	if len(sortText) > 0 && ((sortText != "ASC") && (sortText != "DESC")) {
		return errors.New("Invalid sort parameter. - " + params.Sort)
	}

	// Set order and get query object
	if (len(params.Order) > 0) && (len(params.Sort) > 0) {
		query = t.Order(params.Order + " " + params.Sort)
	} else if len(params.Order) > 0 {
		query = t.Order(params.Order + " ASC")
	} else {
		query = t.Order("id ASC")
	}

	// Are we debugging this?
	if params.Debug {
		query = query.Debug()
	}

	// Offset
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Limit
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	// Add in user id (almost every table has this column)
	if params.UserId > 0 {
		query = query.Where("user_id = ?", params.UserId)
	}

	// Search a particular column
	if (len(params.SearchCol) > 0) && (len(params.SearchTerm) > 0) {
		query = query.Where(params.SearchCol+" LIKE ?", "%"+params.SearchTerm+"%")
	}

	// Run query.
	if err := query.Find(model).Error; err != nil {
		return err
	}

	// If we made it this far no errors.
	return nil
}

/* End File */
