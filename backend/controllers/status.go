//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"app.options.cafe/library/cache"
	"app.options.cafe/library/market"
	"github.com/gin-gonic/gin"
)

//
// Return the market status from tradier
//
func (t *Controller) GetMarketStatus(c *gin.Context) {

	// Get value from our cache
	result := market.MarketStatus{}
	_, err := cache.Get("oc-market-status", &result)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	c.JSON(200, result)
}

/* End File */
