//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
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

	// Set start date
	start, _ := time.Parse("2006-01-02", c.Query("start"))

	if err != nil {
		t.RespondError(c, err, "Unable to parse the start date.")
		return
	}

	// Set end date
	end, _ := time.Parse("2006-01-02", c.Query("end"))

	if err != nil {
		t.RespondError(c, err, "Unable to parse the end date.")
		return
	}

	// Validate the interval
	if !IsValidInterval(c.Query("interval")) {
		t.RespondError(c, err, "Unable to parse the interval.")
		return
	}

	// Make API call to broker.
	result, err := broker.GetHistoricalQuotes(c.Query("symbol"), start, end, c.Query("interval"))

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Return happy JSON
	c.JSON(200, result)
}

// --------------- Helper Functions --------------- //

//
// Test to see if the interval passed in is valid.
//
func IsValidInterval(interval string) bool {
	switch interval {
	case
		"daily",
		"weekly",
		"monthly":
		return true
	}
	return false
}

/* End File */
