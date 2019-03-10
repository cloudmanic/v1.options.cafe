//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/reports"
)

//
// ReportsGetProfitLoss - Return profit and losses
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

	// Figure out if this is Cumulative
	cum := false

	if c.Query("cumulative") == "true" {
		cum = true
	}

	// Get list of profits
	profits := reports.GetProfitLoss(t.DB, brokerAccount, reports.ProfitLossParams{
		StartDate:  helpers.ParseDateNoError(c.Query("start")),
		EndDate:    helpers.ParseDateNoError(c.Query("end")),
		GroupBy:    c.Query("group"),
		Sort:       c.Query("sort"),
		Cumulative: cum,
	})

	// Return happy JSON
	c.JSON(200, profits)
}

//
// ReportsGetTradeGroupYears - Return a list of years that have trade groups
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
// ReportsGetAccountYearlySummary - Get a yearly summary based on account, year
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

//
// ReportsGetAccountReturns - return an array of account returns.
//
func (t *Controller) ReportsGetAccountReturns(c *gin.Context) {
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

	// Get account returns
	r := reports.GetAccountReturns(t.DB, brokerAccount, reports.BalancesParams{
		StartDate: helpers.ParseDateNoError(c.Query("start")),
		EndDate:   helpers.ParseDateNoError(c.Query("end")),
	})

	// Return happy JSON
	c.JSON(200, r)
}

/* End File */
