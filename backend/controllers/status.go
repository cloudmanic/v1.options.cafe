//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/state"
	"github.com/gin-gonic/gin"
)

//
// Return the market status from tradier
//
func (t *Controller) GetMarketStatus(c *gin.Context) {
	c.JSON(200, state.GetMarketStatus())
}

/* End File */
