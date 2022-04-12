//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"app.options.cafe/backtesting"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
)

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

	// Store the backtest
	nBt, err := t.DB.CreateBacktest(o)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Setup a new backtesting & run it.
	bt := backtesting.New(t.DB, int(userId), o.Benchmark)
	go bt.DoBacktestDays(&o)

	// Return happy JSON
	t.RespondJSON(c, http.StatusCreated, gin.H{"id": nBt.Id})
}
