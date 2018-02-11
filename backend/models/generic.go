//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type QueryParam struct {
	UserId           uint
	AccountId        uint
	Limit            uint
	Offset           uint
	Page             string // string because it comes in from the url most likely
	Order            string
	Sort             string
	SearchCols       []string
	SearchTerm       string
	Debug            bool
	Wheres           []KeyValue
	PreLoads         []string
	AllowedOrderCols []string
}

type KeyValue struct {
	Key   string
	Value string
}

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

	// If we passed in a page we figure out the offset from the page.
	if len(params.Page) > 0 {
		page, err := strconv.Atoi(params.Page)

		if err != nil {
			return err
		}

		if (page > 0) && (params.Limit > 0) {
			params.Offset = (uint(page) * params.Limit) - params.Limit
		}
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

	// Add preloads
	for _, row := range params.PreLoads {
		query = query.Preload(row)
	}

	// Add in Where clauses
	for _, row := range params.Wheres {

		if (len(row.Value) > 0) && (len(row.Key) > 0) {
			query = query.Where(row.Key+" = ?", row.Value)
		}

	}

	// Search a particular column
	if (len(params.SearchTerm) > 0) && (len(params.SearchCols) > 0) {
		var likes []string
		var terms []interface{}

		for _, row := range params.SearchCols {
			str := row + " LIKE ?"
			likes = append(likes, str)
			terms = append(terms, "%"+params.SearchTerm+"%")
		}

		// Built where query.
		query = query.Where(strings.Join(likes, " OR "), terms...)
	}

	// Run query.
	if err := query.Find(model).Error; err != nil {
		return err
	}

	// If we made it this far no errors.
	return nil
}

/* End File */
