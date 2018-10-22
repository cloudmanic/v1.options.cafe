//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var (
	NotificationsAllowedOrderCols = []string{"id", "status", "channel", "uri", "sent_tile", "expires"}
)

//
// Update a notifications.
//
func (t *Controller) UpdateNotifications(c *gin.Context) {

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

	// Get the notification in question.
	result := models.Notification{}

	// Run the query
	err = t.DB.Query(&result, models.QueryParam{
		UserId: userId,
		Debug:  true,
		Wheres: []models.KeyValue{
			{Key: "id", ValueInt: int(id)},
		},
	})

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the new status.
	status := gjson.Get(string(body), "status").String()

	// Only status change we can make is to "seen"
	if status != "seen" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status change."})
		return
	}

	// Update status.
	result.Status = status
	t.DB.New().Save(&result)

	spew.Dump(result)

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

//
// Get notifications.
//
func (t *Controller) GetNotifications(c *gin.Context) {

	// Get / Set typical query parms
	page, limit, _ := GetSetPagingParms(c)

	results := []models.Notification{}

	// Run the query
	err := t.DB.Query(&results, models.QueryParam{
		UserId:           c.MustGet("userId").(uint),
		Order:            c.Query("order"),
		Sort:             c.Query("sort"),
		Limit:            limit,
		Page:             page,
		Debug:            false,
		AllowedOrderCols: NotificationsAllowedOrderCols,
		Wheres: []models.KeyValue{
			{Key: "uri", Value: c.Query("uri")},
			{Key: "status", Value: c.Query("status")},
			{Key: "channel", Value: c.Query("channel")},
		},
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)

}

//
// Add a notification channel for this user.
//
func (t *Controller) CreateNotifyChannel(c *gin.Context) {

	// Setup NotifyChannel obj
	o := models.NotifyChannel{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Make sure the UserId is correct.
	o.UserId = c.MustGet("userId").(uint)

	// See if we already have this entry by channelID
	s := models.NotifyChannel{}
	t.DB.New().Where("user_id = ? AND type = ? AND channel_id = ?", o.UserId, o.Type, o.ChannelId).Find(&s)

	// If we already have an obj return that one.
	if s.Id > 0 {
		t.RespondCreated(c, s, nil)
	} else {
		// Create Screen
		err := t.DB.CreateNewRecord(&o, models.InsertParam{})
		t.RespondCreated(c, o, err)
	}
}

/* End File */
