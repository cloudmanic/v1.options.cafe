//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"app.options.cafe/brokers/tradier"
	"app.options.cafe/library/cache"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
)

// Supported screener keys
var screenerItemKeys = []string{
	"open-debit",
	"open-credit",
	"spread-width",
	"days-to-expire",
	"short-strike-percent-away",
	"put-leg-width",
	"call-leg-width",
	"put-leg-percent-away",
	"call-leg-percent-away",
}

// Supported screener operator
var screenerItemOperators = []string{
	"=",
	">",
	"<",
}

//
// Delete a screener
//
func (t *Controller) DeleteScreener(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the screener by id.
	orgObj, err := t.DB.GetScreenerByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Delete items.
	t.DB.New().Where("screener_id = ?", orgObj.Id).Delete(models.ScreenerItem{})

	// Delete record
	t.DB.New().Delete(&orgObj)

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

//
// Update a screener
//
func (t *Controller) UpdateScreener(c *gin.Context) {

	// Setup Screener obj
	o := models.Screener{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Make sure the UserId is correct.
	o.UserId = c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the screener by id.
	orgObj, err := t.DB.GetScreenerByIdAndUserId(uint(id), o.UserId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Set id.
	o.Id = orgObj.Id

	// Set the create date.
	o.CreatedAt = orgObj.CreatedAt

	// This is hacky but allows us to reset items
	t.DB.New().Where("screener_id = ?", orgObj.Id).Delete(models.ScreenerItem{})

	// Add in UserId to items & Some extra validation for ScreenerItems
	if !addUserIdAndValidateScreenerItems(&o, c) {
		return
	}

	// Update Screen
	t.DB.New().Save(&o)

	// Clear any cache
	cache.Delete("oc-screener-result-" + strconv.Itoa(int(orgObj.Id)))

	// Return success.
	c.JSON(http.StatusNoContent, nil)
}

//
// Create new screener
//
func (t *Controller) CreateScreener(c *gin.Context) {

	// Setup Screener obj
	o := models.Screener{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Make sure the UserId is correct.
	o.UserId = c.MustGet("userId").(uint)

	// Add in UserId to items & Some extra validation for ScreenerItems
	if !addUserIdAndValidateScreenerItems(&o, c) {
		return
	}

	// Create Screen
	err := t.DB.CreateNewRecord(&o, models.InsertParam{})
	t.RespondCreated(c, o, err)
}

//
// Get a screeners
//
func (t *Controller) GetScreeners(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the screener by id.
	screeners, err := t.DB.GetScreenersByUserId(userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, screeners)
}

//
// Get a screener by id.
//
func (t *Controller) GetScreener(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the screener by id.
	screener, err := t.DB.GetScreenerByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, screener)
}

//
// Get screen results
//
func (t *Controller) GetScreenerResults(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the screener by id.
	screen, err := t.DB.GetScreenerByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Get trade groups so we know which positions we already have on.
	tradeGroups, _, _ := t.DB.GetTradeGroups(models.QueryParam{
		UserId:     userId,
		Order:      "open_date",
		Sort:       "asc",
		Limit:      1000,
		Page:       1,
		Debug:      false,
		PreLoads:   []string{"Positions"},
		SearchTerm: "",
		SearchCols: []string{"id", "name", "open_date", "status", "type", "note"},
		Wheres: []models.KeyValue{
			{Key: "status", Value: "Open"},
			{Key: "broker_account_id", Value: c.Query("broker_account_id")},
		},
	})

	// Loop through and just get the symbols
	positions := []uint{}

	for _, row := range tradeGroups {
		for _, row2 := range row.Positions {
			positions = append(positions, row2.Symbol.Id)
		}
	}

	// See if we have this result in the cache.
	// We keep the cache up to date via a feed loop started from main.go
	cachedResult := []screener.Result{}

	found, _ := cache.Get("oc-screener-result-"+strconv.Itoa(int(screen.Id)), &cachedResult)

	// Return happy JSON
	if found {
		c.JSON(200, returnPositionOns(cachedResult, positions))
		return
	}

	// Setup the broker
	broker := tradier.Api{DB: nil, ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}

	// New screener instance
	s := screener.NewScreen(t.DB, &broker)

	// Run the screen based from our function map
	result, err := s.ScreenFuncs[screen.Strategy](screen)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, returnPositionOns(result, positions))
}

//
// Get screen from filters
//
func (t *Controller) GetScreenerResultsFromFilters(c *gin.Context) {

	// Setup Screener obj
	screen := models.Screener{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &screen) != nil {
		return
	}

	// Take md5 of the status
	hash, err := structhash.Hash(screen, 1)

	if t.RespondError(c, err, httpNoRecordFound) {
		services.Info(err)
		return
	}

	// See if we have this result in the cache.
	cachedResult := []screener.Result{}

	found, _ := cache.Get("oc-screener-result-"+hash, &cachedResult)

	// Return happy JSON
	if found {
		c.JSON(200, cachedResult)
		return
	}

	// Setup the broker
	broker := tradier.Api{DB: nil, ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}

	// New screener instance
	s := screener.NewScreen(t.DB, &broker)

	// Run the screen based from our function map
	result, err := s.ScreenFuncs[screen.Strategy](screen)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Store result in cache.
	cache.SetExpire("oc-screener-result-"+hash, (time.Minute * 5), result)

	// Return happy JSON
	c.JSON(200, result)
}

// ----------- Helper Function ------------- //

//
// Add in UserId to items & Some extra validation for ScreenerItems
//
func addUserIdAndValidateScreenerItems(o *models.Screener, c *gin.Context) bool {

	for key, row := range o.Items {
		o.Items[key].UserId = o.UserId

		// Validation - screenerItemKeys
		found, _ := helpers.InArray(row.Key, screenerItemKeys)

		if !found {
			m := make(map[string]string)
			m["items"] = "Unknown Key - " + row.Key + "."
			c.JSON(http.StatusBadRequest, gin.H{"errors": m})
			return false
		}

		// Validation - screenerItemOperators
		found2, _ := helpers.InArray(row.Operator, screenerItemOperators)

		if !found2 {
			m := make(map[string]string)
			m["items"] = "Unknown operator - " + row.Operator + "."
			c.JSON(http.StatusBadRequest, gin.H{"errors": m})
			return false
		}

	}

	// Return happy.
	return true
}

//
// returnPositionOns will set the boolean of symbols we already have on.
//
func returnPositionOns(results []screener.Result, positions []uint) []screener.Result {
	// Match up any found positions
	for key, row := range results {
		for key2, row2 := range row.Legs {
			if findPositionFromTradeGroups(positions, row2.Id) {
				results[key].Legs[key2].PositionOn = true
			}
		}
	}

	return results
}

//
// Search for positions in our trade groups.
//
func findPositionFromTradeGroups(positions []uint, symbolId uint) bool {
	for _, row := range positions {
		if row == symbolId {
			return true
		}
	}

	return false
}

/* End File */
