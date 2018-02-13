//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Return groups in our database.
//
func (t *Controller) GetBrokers(c *gin.Context) {

	// // Run the query
	// results, err := t.DB.GetTradeGroups(models.QueryParam{
	// 	UserId:           c.MustGet("userId").(uint),
	// 	Order:            c.Query("order"),
	// 	Sort:             c.Query("sort"),
	// 	Limit:            defaultMysqlLimit,
	// 	Page:             c.Query("page"),
	// 	Debug:            false,
	// 	PreLoads:         []string{"Positions"},
	// 	SearchTerm:       c.Query("search"),
	// 	SearchCols:       []string{"id", "name", "open_date", "status", "type", "note"},
	// 	AllowedOrderCols: TradeGroupAllowedOrderCols,
	// 	Wheres:           []models.KeyValue{{Key: "broker_account_id", Value: c.Query("broker_account_id")}},
	// })

	var results = []models.Broker{}

	// Run the query
	err := t.DB.Query(&results, models.QueryParam{
		UserId:   c.MustGet("userId").(uint),
		Limit:    defaultMysqlLimit,
		PreLoads: []string{"BrokerAccounts"},
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
