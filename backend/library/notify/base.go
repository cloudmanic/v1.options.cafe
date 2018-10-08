//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package notify

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/notify/sms_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify/web_push"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

var channels []string = []string{"web-push", "sms-push"}

type NotifyRequest struct {
	UserId   uint
	Uri      string
	UriRefId uint
	ShortMsg string
	LongMsg  string
	Status   string
	Title    string
}

//
// Push to all channels. Mata is used to pass in any general infro
//
func Push(db models.Datastore, request NotifyRequest) {

	// Loop through the different channels and push
	for _, row := range channels {
		go PushChannel(db, row, request)
	}

}

//
// Push down one channel
//
func PushChannel(db models.Datastore, channel string, request NotifyRequest) {

	// Set some helpful times
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Set title
	title := "Options Cafe Trading"

	if len(request.Title) > 0 {
		title = request.Title
	}

	// See if we have already sent this notification
	o := models.Notification{}
	db.New().Where("sent_time > ? AND channel = ? AND user_id = ? AND uri = ? AND uri_ref_id = ?", dayStart, channel, request.UserId, request.Uri, request.UriRefId).Find(&o)

	// We already sent this notice
	if o.Id > 0 {
		return
	}
	// Store this notification
	ob := models.Notification{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserId:       request.UserId,
		Status:       "pending",
		Channel:      channel,
		Uri:          request.Uri,
		UriRefId:     request.UriRefId,
		Title:        title,
		ShortMessage: request.ShortMsg,
		LongMessage:  request.LongMsg,
		SentTime:     time.Now(),
		Expires:      time.Now(),
	}

	// Store in DB
	db.New().Save(&ob)

	// Special case for market status.
	if (request.Uri == "market-status-open") || (request.Uri == "market-status-closed") {

		switch channel {

		case "web-push":
			DoMarketStatusChangeWebPush(db, title, request.Status, request.ShortMsg)

		case "sms-push":
			DoMarketStatusChangeSmsPush(db, request.Status, request.ShortMsg)

		}

	} else {

		// TODO: We need to look up in settings if the user wants this notification. On Market status we check within that function.

		switch channel {

		case "web-push":
			DoWebPushForUser(db, request.UserId, title, request.Status, request.ShortMsg)

		case "sms-push":
			DoSmsPushForUser(db, request.UserId, request.Status, request.ShortMsg)

		}

	}

	// Update Notification as sent.
	ob.Status = "sent"
	db.New().Save(&ob)

}

//
// Send SMS Push for a particular user.
//
func DoSmsPushForUser(db models.Datastore, userId uint, status string, content string) {

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ? AND user_id = ?", "sms-push", userId).Find(&nc)

	// Loop through and send message
	for _, row := range nc {
		sms_push.Push(row.ChannelId, content)
	}

}

//
// Send Web Push for a particular user.
//
func DoWebPushForUser(db models.Datastore, userId uint, title string, status string, content string) {

	deviceIds := []string{}

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ? AND user_id = ?", "web-push", userId).Find(&nc)

	for _, row := range nc {
		deviceIds = append(deviceIds, row.ChannelId)
	}

	// Send message
	if len(deviceIds) > 0 {
		web_push.Push(deviceIds, title, content)
	}

}

//
// Do market status. Since this goes out to everyone generically we
// handle it differently.
//
func DoMarketStatusChangeSmsPush(db models.Datastore, status string, content string) {

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "closed") {
		return
	}

	// TODO: We need to look up in settings if the user wants this notification.

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "sms-push").Find(&nc)

	// Loop through and send message
	for _, row := range nc {
		sms_push.Push(row.ChannelId, content)
	}

}

//
// Do market status. Since this goes out to everyone generically we
// handle it differently.
//
func DoMarketStatusChangeWebPush(db models.Datastore, title string, status string, content string) {

	deviceIds := []string{}

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "closed") {
		return
	}

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "web-push").Find(&nc)

	for _, row := range nc {

		// TODO: We need to look up in settings if the user wants this notification.

		deviceIds = append(deviceIds, row.ChannelId)
	}

	// Send message
	if len(deviceIds) > 0 {
		web_push.Push(deviceIds, title, content)
	}

}

/* End File */
