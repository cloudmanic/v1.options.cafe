//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"flag"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
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
// Get a brokers active orders.
//
func (t *Controller) GetBrokerActiveOrders(c *gin.Context) {

	// Build cache key
	key := "oc-orders-active-" + strconv.Itoa(int(c.MustGet("userId").(uint))) + "-" + string(c.Param("id"))

	// Get a value we know we do not have
	result := []types.Order{}
	_, err := cache.Get(key, &result)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy.
	c.JSON(200, result)
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

	// // Figure out which broker connection to setup.
	// switch broker.Name {

	// case "Tradier":
	// 	brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: false}

	// case "Tradier Sandbox":
	// 	brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: true}

	// default:
	// 	services.Critical("Order: Unknown Broker : " + broker.Name)

	// }

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
