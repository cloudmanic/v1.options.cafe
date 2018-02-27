//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"flag"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Get time sales data
//
// TODO: pull the access token from the user's account.
//
func (t *Controller) GetHistoricalQuotes(c *gin.Context) {

	var apiKey string = ""
	var brokers = []models.Broker{}

	// Run the query to get brokers
	err := t.DB.Query(&brokers, models.QueryParam{
		UserId: c.MustGet("userId").(uint),
		Wheres: []models.KeyValue{
			{Key: "name", Value: "Tradier"},
		},
	})

	// TODO: For now we only support Tradier but as we open up to new brokers we will have to support more.
	if flag.Lookup("test.v") == nil {
		for _, row := range brokers {

			// Decrypt the access token
			_apiKey, err := helpers.Decrypt(row.AccessToken)

			if err != nil {
				t.RespondError(c, err, httpGenericErrMsg)
				return
			}

			apiKey = _apiKey
		}
	}

	// Setup the broker
	broker := tradier.Api{
		DB:     t.DB,
		ApiKey: apiKey,
	}

	// Validate the interval
	if !IsValidInterval(c.Query("interval")) {
		t.RespondError(c, errors.New("Unable to parse the interval. - "+c.Query("interval")), "Unable to parse the interval.")
		return
	}

	// Set start date
	start, _ := time.Parse("2006-01-02 15:04", c.Query("start")+" 09:30")

	if err != nil {
		t.RespondError(c, err, "Unable to parse the start date.")
		return
	}

	// Set end date
	end, _ := time.Parse("2006-01-02 15:04", c.Query("end")+" 16:00")

	if err != nil {
		t.RespondError(c, err, "Unable to parse the end date.")
		return
	}

	// Store in active symbols
	t.DB.CreateActiveSymbol(c.MustGet("userId").(uint), c.Query("symbol"))

	// Make API call to broker.
	if IsHistorical(c.Query("interval")) {

		// Make call to broker.
		result, err := broker.GetHistoricalQuotes(c.Query("symbol"), start, end, c.Query("interval"))

		if err != nil {
			t.RespondError(c, err, httpGenericErrMsg)
			return
		}

		// Return happy JSON
		c.JSON(200, result)

	} else {

		// Make call to broker.
		result, err := broker.GetTimeSalesQuotes(c.Query("symbol"), start, end, c.Query("interval"))

		if err != nil {
			t.RespondError(c, err, httpGenericErrMsg)
			return
		}

		// Return happy JSON
		c.JSON(200, result)

	}
}

// --------------- Helper Functions --------------- //

//
// Test to see if this is a historical quote or a time sales quote
//
func IsHistorical(interval string) bool {
	switch interval {
	case
		"daily",
		"weekly",
		"monthly":
		return true
	}
	return false
}

//
// Test to see if the interval passed in is valid.
//
func IsValidInterval(interval string) bool {
	switch interval {
	case
		"1min",
		"5min",
		"15min",
		"daily",
		"weekly",
		"monthly":
		return true
	}
	return false
}

/* End File */
