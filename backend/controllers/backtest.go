//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"strconv"

	"app.options.cafe/library/queue"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
)

//
// GetBacktests will return all backtests to a user.
//
func (t *Controller) GetBacktests(c *gin.Context) {
	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the backtests by user id.
	bts, err := t.DB.BacktestsGetByUserId(userId, true)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, bts)
}

//
// GetBacktest will return one backtest result
//
func (t *Controller) GetBacktest(c *gin.Context) {
	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the backtests by user id.
	bt, err := t.DB.BacktestGetByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, bt)
}

//
// CreateBacktest will create and start a new backtest.
//
func (t *Controller) CreateBacktest(c *gin.Context) {
	// Setup Backtest obj
	o := models.Backtest{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Get the user id.
	userId := c.MustGet("userId").(uint)
	o.UserId = userId
	o.Screen.UserId = userId
	o.EndingBalance = o.StartingBalance // They will be the same at the start with no trades

	// Store the backtest
	nBt, err := t.DB.CreateBacktest(o)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Send message to start backtest
	queue.Write("oc-job", `{"action":"backtest-run-days","user_id":`+strconv.Itoa(int(userId))+`,"backtest_id":`+strconv.Itoa(int(nBt.Id))+`}`)

	// Return happy JSON
	t.RespondJSON(c, http.StatusCreated, gin.H{"id": nBt.Id})
}
