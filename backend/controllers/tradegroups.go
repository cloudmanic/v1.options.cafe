//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

var (
	TradeGroupAllowedOrderCols = []string{"id", "open_date", "closed_date", "profit"}
)

//
// Return groups in our database.
//
func (t *Controller) GetTradeGroups(c *gin.Context) {

	// Run the query
	results, err := t.DB.GetTradeGroups(models.QueryParam{
		UserId:           c.MustGet("userId").(uint),
		Order:            c.Query("order"),
		Sort:             c.Query("sort"),
		Limit:            defaultMysqlLimit,
		Page:             c.Query("page"),
		Debug:            false,
		PreLoads:         []string{"Positions"},
		SearchTerm:       c.Query("search"),
		SearchCols:       []string{"id", "name", "open_date", "status", "type", "note"},
		AllowedOrderCols: TradeGroupAllowedOrderCols,
		Wheres: []models.KeyValue{
			{Key: "type", Value: c.Query("type")},
			{Key: "status", Value: c.Query("status")},
			{Key: "broker_account_id", Value: c.Query("broker_account_id")},
		},
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
