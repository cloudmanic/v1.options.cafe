//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// Return a full option chain based on symbol and expire date.
//
func (t *Controller) GetOptionsChainByExpiration(c *gin.Context) {

	// Get access token
	apiKey, err := t.GetTradierAccessToken(c)

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Setup the broker
	broker := tradier.Api{
		DB:     t.DB,
		ApiKey: apiKey,
	}

	// Get chain from tradier.
	chain, err := broker.GetOptionsChainByExpiration(c.Param("symb"), c.Param("expire"))

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Return happy JSON
	c.JSON(200, chain)
}

//
// Pass in a symbol and expire date to get a list of strikes.
// For now we just use our admin Tradier API key. Some day we
// should redo this to use the users API key.
//
func (t *Controller) GetOptionsStikesBySymbolExpiration(c *gin.Context) {

	var result []int64

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://api.tradier.com/v1/markets/options/strikes?symbol="+c.Param("symb")+"&expiration="+c.Param("expire"), nil)

	// Headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"))
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	res, err := client.Do(req)

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Close Body
	defer res.Body.Close()

	// Read Response Body
	json, _ := ioutil.ReadAll(res.Body)

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		t.RespondError(c, errors.New("Failed response from Tradier."), httpGenericErrMsg)
		return
	}

	// Loop through the strikes
	s := gjson.Get(string(json), "strikes.strike")
	for _, row := range s.Array() {
		result = append(result, row.Int())
	}

	// Return happy JSON
	c.JSON(200, result)
}

//
// Pass in a symbol and return a list of all expirations for
// this options symbol. Just to keep things simple we us our
// admin account with Tradier to get this information. At this point
// no need to support this with every broker and such.
//
func (t *Controller) GetOptionsExpirations(c *gin.Context) {

	var result []string

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://api.tradier.com/v1/markets/options/expirations?symbol="+c.Param("symb"), nil)

	// Headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"))
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	res, err := client.Do(req)

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Close Body
	defer res.Body.Close()

	// Read Response Body
	json, _ := ioutil.ReadAll(res.Body)

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		t.RespondError(c, errors.New("Failed response from Tradier."), httpGenericErrMsg)
		return
	}

	// Loop through the dates
	dates := gjson.Get(string(json), "expirations.date")
	for _, row := range dates.Array() {
		result = append(result, row.String())
	}

	// Return happy JSON
	c.JSON(200, result)
}

//
// Get time sales data
//
// TODO: pull the access token from the user's account.
//
func (t *Controller) GetHistoricalQuotes(c *gin.Context) {

	// Get access token
	apiKey, err := t.GetTradierAccessToken(c)

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
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
	start, err := time.Parse("2006-01-02 15:04", c.Query("start")+" 09:30")

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
// Get access tradier access token.
//
func (t *Controller) GetTradierAccessToken(c *gin.Context) (string, error) {

	var apiKey string = ""
	var brokers = []models.Broker{}

	// Run the query to get brokers
	err := t.DB.Query(&brokers, models.QueryParam{
		UserId: c.MustGet("userId").(uint),
		Wheres: []models.KeyValue{
			{Key: "name", Value: "Tradier"},
		},
	})

	if err != nil {
		return "", err
	}

	// TODO: For now we only support Tradier but as we open up to new brokers we will have to support more.
	if flag.Lookup("test.v") == nil {
		for _, row := range brokers {

			// Decrypt the access token
			_apiKey, err := helpers.Decrypt(row.AccessToken)

			if err != nil {
				return "", err
			}

			apiKey = _apiKey
		}
	}

	// Return Happy
	return apiKey, nil
}

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
