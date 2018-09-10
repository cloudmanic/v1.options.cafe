//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// Get a yearly summary based on account, year
//
func (t *Controller) ReportsGetAccountYearlySummary(c *gin.Context) {

	// Make sure the UserId is correct.
	//userId := c.MustGet("userId").(uint)

	// // Set as int
	// id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	// if t.RespondError(c, err, httpGenericErrMsg) {
	// 	return
	// }

	// // Get the screener by id.
	// orgObj, err := t.DB.GetScreenerByIdAndUserId(uint(id), userId)

	// if t.RespondError(c, err, httpNoRecordFound) {
	// 	return
	// }

	// // Delete items.
	// t.DB.New().Where("screener_id = ?", orgObj.Id).Delete(models.ScreenerItem{})

	// // Delete record
	// t.DB.New().Delete(&orgObj)

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

/* End File */
