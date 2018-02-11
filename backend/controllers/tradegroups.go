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
	AllowedOrderCols = []string{"id", "open_date", "closed_date", "profit"}
)

//
// Return groups in our database.
//
func (t *Controller) GetTradeGroups(c *gin.Context) {

	// Place to store the results.
	var results = []models.TradeGroup{}

	// Run the query
	err := t.DB.Query(&results, models.QueryParam{
		UserId:           c.MustGet("userId").(uint),
		Order:            c.Query("order"),
		Sort:             c.Query("sort"),
		Limit:            defaultMysqlLimit,
		Offset:           0,
		Debug:            true,
		AllowedOrderCols: AllowedOrderCols,
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
