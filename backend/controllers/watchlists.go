//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// Return watchlists in our database.
//
func (t *Controller) GetWatchlists(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the watchlists
	wLists, err := t.DB.GetWatchlistsByUserId(userId)

	if t.RespondError(c, err, "No access to this watchlist resource.") {
		return
	}

	// Return happy JSON
	c.JSON(200, wLists)
}

//
// Return watchlist in our database.
//
func (t *Controller) GetWatchlist(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the watchlist by id.
	wList, err := t.DB.GetWatchlistsByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	c.JSON(200, wList)
}

//
// Watchlist - Create
//
// curl -H "Content-Type: application/json" -X POST -d '{"name":"Super Cool Watchlist"}' -H "Authorization: Bearer XXXXXX" http://localhost:7080/api/v1/watchlists
//
func (t *Controller) CreateWatchlist(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	name := gjson.Get(string(body), "name").String()

	// Validate name
	if len(name) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name field can not be empty."})
		return
	}

	// Validate if this watchlist is already in the db.
	if !t.DB.New().Where("user_id = ? AND name = ?", userId, name).First(&models.Watchlist{}).RecordNotFound() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A watchlist already exists by this name."})
		return
	}

	// Get the watchlists
	wLists, err := t.DB.CreateWatchlist(userId, name)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	t.RespondJSON(c, http.StatusOK, wLists)
}

//
// Update a watchlist.
//
func (t *Controller) UpdateWatchlist(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the watchlist by id.
	_, err = t.DB.GetWatchlistsByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Get the new watchlist name.
	name := gjson.Get(string(body), "name").String()

	// Validate name
	if len(strings.Trim(name, " ")) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name field can not be empty."})
		return
	}

	// Validate if this watchlist is already in the db.
	if !t.DB.New().Where("user_id = ? AND name = ?", userId, strings.Trim(name, " ")).First(&models.Watchlist{}).RecordNotFound() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A watchlist already exists by this name."})
		return
	}

	// Update the name of the watchlist.
	err = t.DB.WatchlistUpdate(uint(id), name)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})

}

//
// Delete a watchlist.
//
func (t *Controller) DeleteWatchlist(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the watchlist by id.
	_, err = t.DB.GetWatchlistsByIdAndUserId(uint(id), userId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Now delete the watchlist.
	err = t.DB.WatchlistDeleteById(uint(id))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

//
// Add a symbol to a watch list. Inserts into the WatchlistSymbol model.
//
func (t *Controller) WatchlistAddSymbol(c *gin.Context) {

	// Setup WatchlistSymbol obj
	o := models.WatchlistSymbol{
		UserId:      c.MustGet("userId").(uint),
		WatchlistId: helpers.StringToUint(c.Param("id")),
	}

	// Validate if a user has access to this watchlist.
	if !t.ValidateWatchlistUserAccess(o.UserId, o.WatchlistId) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access to this watchlist resource."})
		return
	}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Prepends symbol on to the watchlist.
	err := t.DB.PrependWatchlistSymbol(&o)
	t.RespondCreated(c, o, err)
}

//
// Reorder a watch list.
//
func (t *Controller) WatchlistReorder(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Validate if a user has access to this watchlist.
	if !t.ValidateWatchlistUserAccess(userId, uint(id)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access to this watchlist resource."})
		return
	}

	// Ids was into the model.
	var ids []int

	// Parse json body
	jsonBody, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	result := gjson.Get(string(jsonBody), "ids")

	for _, id := range result.Array() {
		ids = append(ids, int(id.Int()))
	}

	// Send the reorder into our model
	err = t.DB.WatchlistReorder(uint(id), ids)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

//
// Delete a symbol from a watchlist.
func (t *Controller) WatchlistDeleteSymbol(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Set as int - symb
	symbId, err := strconv.ParseInt(c.Param("symb"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Validate if a user has access to this watchlist.
	if !t.ValidateWatchlistUserAccess(userId, uint(id)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access to this watchlist resource."})
		return
	}

	// Delete symbol from watchlist.
	err = t.DB.WatchlistRemoveSymbol(uint(id), uint(symbId))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

// ----------------- Helper Functions -------------------- //

//
// Validate a user has access to a passed in watchlist.
//
func (t *Controller) ValidateWatchlistUserAccess(u uint, w uint) bool {
	_, err := t.DB.GetWatchlistsByIdAndUserId(w, u)

	if err != nil {
		return false
	}

	return true
}

/* End File */
