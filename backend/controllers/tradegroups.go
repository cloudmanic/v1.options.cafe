//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"
)

//
// Return groups in our database.
//
func (t *Controller) GetTradeGroups(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get basic query parms
	parms := t.GetBasicQueryValues(c)

	// Get the watchlists
	groups, err := t.DB.GetTradeGroupsByUserId(userId, parms.FullOrder)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, groups)
}

/* End File */
