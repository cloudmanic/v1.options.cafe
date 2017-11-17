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

	"app.options.cafe/backend/library/services"
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

	if t.RespondError(c, err, httpGenericErrMsg) {
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

	// // // Get the user id.
	userId := c.MustGet("userId").(uint)

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	name := gjson.Get(string(body), "name").String()

	// Get the watchlists
	wLists, err := t.DB.CreateWatchlist(userId, name)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	json := t.RespondJSON(c, http.StatusOK, wLists)

	// Send new watchlist up websocket.
	t.SendWatchlistUpWS(userId, json)
}

//
// Send watchlist up websocket.
//
func (t *Controller) SendWatchlistUpWS(userId uint, json string) {

	// Build JSON we send
	jsonSend, err := t.WsSendJsonBuild("watchlists", json)

	if err != nil {
		services.Error(err, "GetWatchlists() WsSendJsonBuild (#1)")
		return
	}

	// Send new watchlist through the websocket.
	t.WsWriteChan <- SendStruct{UserId: userId, Body: jsonSend}
}

/* End File */
