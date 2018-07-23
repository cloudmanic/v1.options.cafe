//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
	"github.com/gin-gonic/gin"
)

// Supported screener keys
var screenerItemKeys = []string{
	"min-credit",
	"spread-width",
	"max-days-to-expire",
	"min-days-to-expire",
	"short-strike-percent-away",
}

// Supported screener operator
var screenerItemOperators = []string{
	"=",
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

	// See if we have this result in the cache.
	cachedResult := []screener.Result{}

	found, _ := cache.Get("oc-screener-result-"+strconv.Itoa(int(screen.Id)), &cachedResult)

	// Return happy JSON
	if found {
		c.JSON(200, cachedResult)
		return
	}

	// Run back test
	result, err := screener.RunPutCreditSpread(screen, t.DB)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Store result in cache.
	cache.SetExpire("oc-screener-result-"+strconv.Itoa(int(screen.Id)), (time.Minute * 5), result)

	// Return happy JSON
	c.JSON(200, result)
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

/* End File */
