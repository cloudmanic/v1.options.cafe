//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

//
// Get a screeners
//
func (t *Controller) GetScreeners(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the screener by id.
	screeners, err := t.DB.GetScreenersByUserId(userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, screeners)
}

//
// Get a screener by id.
//
func (t *Controller) GetScreener(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the screener by id.
	screener, err := t.DB.GetScreenerByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, screener)
}

/* End File */
