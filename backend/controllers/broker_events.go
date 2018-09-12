//
// Date: 9/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

var (
	BrokerEventsAllowedOrderCols = []string{"id", "date", "symbol", "type", "trade_type"}
)

//
// Return a list of broker events.
//
func (t *Controller) GetBrokerEvents(c *gin.Context) {

	// Set as int - brokerAccountId
	brokerAccountId, err := strconv.ParseInt(c.Param("brokerAccount"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get / Set typical query parms
	page, limit, _ := GetSetPagingParms(c)

	// Setup resolt.
	results := []models.BrokerEvent{}

	// Setup the query params
	params := models.QueryParam{
		UserId:           c.MustGet("userId").(uint),
		Order:            c.Query("order"),
		Sort:             c.Query("sort"),
		Limit:            limit,
		Page:             page,
		Debug:            false,
		PreLoads:         []string{},
		SearchTerm:       c.Query("search"),
		SearchCols:       []string{"id", "type", "date", "amount", "symbol", "commission", "description", "price", "quantity", "trade_type"},
		AllowedOrderCols: BrokerEventsAllowedOrderCols,
		Wheres: []models.KeyValue{
			{Key: "broker_account_id", ValueInt: int(brokerAccountId)},
		},
	}

	// Run the query
	noFilterCount, err := t.DB.QueryWithNoFilterCount(&results, params)

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the meta data related to this query.
	meta := t.DB.GetQueryMetaData(len(results), noFilterCount, params)

	// Set some headers for paging.
	t.AddPagingInfoToHeaders(c, meta)

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
