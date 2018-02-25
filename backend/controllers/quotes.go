//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"flag"
	"os"

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
	if flag.Lookup("test.v") != nil {

		// Get API Key if we are in testing mode
		if flag.Lookup("test.v") != nil {
			apiKey = os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")
		}

	} else {

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

	// Make API call to broker.
	result, err := broker.GetBalances()

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Return happy JSON
	c.JSON(200, result)
}

/* End File */
