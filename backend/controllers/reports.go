//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/reports"
	"github.com/gin-gonic/gin"
)

//
// Return profit and losses
//
func (t *Controller) ReportsGetProfitLoss(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Set as int - brokerAccountId
	brokerAccountId, err := strconv.ParseInt(c.Param("brokerAccount"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get broker account
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(uint(brokerAccountId), userId)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get list of profits
	profits := reports.GetProfitLoss(t.DB, brokerAccount, reports.ProfitLossParams{
		StartDate: helpers.ParseDateNoError(c.Query("start")),
		EndDate:   helpers.ParseDateNoError(c.Query("end")),
		GroupBy:   c.Query("group"),
		Sort:      c.Query("sort"),
	})

	// Return happy JSON
	c.JSON(200, profits)
}

//
// Return a list of years that have trade groups
//
func (t *Controller) ReportsGetTradeGroupYears(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Set as int - brokerAccountId
	brokerAccountId, err := strconv.ParseInt(c.Param("brokerAccount"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get broker account
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(uint(brokerAccountId), userId)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get list of years
	years := reports.GetYearsWithTradeGroups(t.DB, brokerAccount)

	// Return happy JSON
	c.JSON(200, years)
}

//
// Get a yearly summary based on account, year
//
func (t *Controller) ReportsGetAccountYearlySummary(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Set as int - Year
	year, err := strconv.ParseInt(c.Param("year"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Set as int - brokerAccountId
	brokerAccountId, err := strconv.ParseInt(c.Param("brokerAccount"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get broker account
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(uint(brokerAccountId), userId)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get summary from database
	summary := reports.GetYearlySummaryByAccountYear(t.DB, brokerAccount, int(year))

	// Return happy JSON
	c.JSON(200, summary)
}

/* End File */
