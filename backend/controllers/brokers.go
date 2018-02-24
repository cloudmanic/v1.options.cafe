//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
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
// Return groups in our database.
//
func (t *Controller) GetBrokers(c *gin.Context) {

	var results = []models.Broker{}

	// Run the query
	err := t.DB.Query(&results, models.QueryParam{
		UserId:   c.MustGet("userId").(uint),
		Limit:    defaultMysqlLimit,
		PreLoads: []string{"BrokerAccounts"},
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)
}

//
// Make an API call to the broker and get balances.
//
func (t *Controller) GetBalances(c *gin.Context) {

	var apiKey string = ""
	var brokers = []models.Broker{}

	// Run the query to get brokers
	err := t.DB.Query(&brokers, models.QueryParam{
		UserId: c.MustGet("userId").(uint),
		Wheres: []models.KeyValue{
			{Key: "name", Value: "Tradier"},
		},
	})

	// Loop through the different brokers- TODO: This only supports one broker. We need to get balance from all brokers and merge data together.
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
