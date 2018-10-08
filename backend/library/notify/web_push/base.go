//
// Date: 2018-10-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-08
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package web_push

import (
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	onesignal "github.com/tbalthazar/onesignal-go"
	"github.com/tidwall/gjson"
)

//
// Push to all web push channels that the users has opted into.
//
func Push(db models.Datastore, userId uint, uri string, uriRefId uint, data_json string) {

	// Set some helpful times
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Get the status & message.
	title := "Options Cafe Trading"
	status := gjson.Get(data_json, "status").String()
	msg := "The market is now " + status

	// See if we have already sent this notification
	o := models.Notification{}
	db.New().Where("sent_time > ? AND status = ? AND channel = ? AND user_id = ? AND uri = ? AND uri_ref_id = ?", dayStart, "sent", "web-push", userId, uri, uriRefId).Find(&o)

	// We already sent this notice
	if o.Id > 0 {
		return
	}

	// Switch based on URI
	switch uri {

	case "market-status-open":
		DoMarketStatusChange(db, title, status, msg)

	case "market-status-closed":
		DoMarketStatusChange(db, title, status, msg)

	}

	// Store this notification
	ob := models.Notification{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserId:       userId,
		Status:       "sent",
		Channel:      "web-push",
		Uri:          uri,
		UriRefId:     uint(0),
		Title:        title,
		ShortMessage: msg,
		LongMessage:  msg,
		SentTime:     time.Now(),
		Expires:      time.Now(),
	}

	// Store in DB
	db.New().Save(&ob)

}

//
// Do market status. Since this goes out to everyone generically we
// handle it differently.
//
func DoMarketStatusChange(db models.Datastore, title string, status string, content string) {

	deviceIds := []string{}

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "closed") {
		return
	}

	// TODO: We need to look up in settings if the user wants this notification.

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "web-push").Find(&nc)

	for _, row := range nc {
		deviceIds = append(deviceIds, row.ChannelId)
	}

	// Send message
	if len(deviceIds) > 0 {
		DoOneSignalWebPushSend(deviceIds, title, content)
	}

}

//
// Send message up to open signal to be pushed out.
//
func DoOneSignalWebPushSend(deviceIds []string, heading string, content string) {

	if len(os.Getenv("ONESIGNAL_APP_ID")) > 0 {

		// Log
		services.Info("Sending web push notification: " + content)

		// One Signal Stuff.
		client := onesignal.NewClient(nil)

		// Create send request.
		notificationReq := &onesignal.NotificationRequest{
			AppID:            os.Getenv("ONESIGNAL_APP_ID"),
			Headings:         map[string]string{"en": heading},
			Contents:         map[string]string{"en": content},
			IsAnyWeb:         true,
			IncludePlayerIDs: deviceIds,
		}

		// Send the request
		_, _, err := client.Notifications.Create(notificationReq)

		if err != nil {
			services.BetterError(err)
		}

	}

}

/* End File */
