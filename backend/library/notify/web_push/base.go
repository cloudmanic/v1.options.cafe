//
// Date: 2018-10-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package web_push

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	onesignal "github.com/tbalthazar/onesignal-go"
	"github.com/tidwall/gjson"
)

//
// Push to all web push channels that the users has opted into.
//
func Push(db models.Datastore, userId uint, uri string, data_json string) {

	switch uri {
	case "market-status":
		DoMarketStatusChange(db, data_json)
	}

}

//
// Do market status. Since this goes out to everyone generically we
// handle it differently.
//
func DoMarketStatusChange(db models.Datastore, data_json string) {

	deviceIds := []string{}

	// Get the status.
	status := gjson.Get(data_json, "status").String()

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "closed") {
		return
	}

	// TODO: We need to look up in settings if the user wants this notification.

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "Web Push").Find(&nc)

	for _, row := range nc {
		deviceIds = append(deviceIds, row.ChannelId)
	}

	// Send message
	if len(deviceIds) > 0 {
		DoOneSignalWebPushSend(deviceIds, "Options Cafe Trading", "The market is now "+status)
	}

}

//
// Send message up to open signal to be pushed out.
//
func DoOneSignalWebPushSend(deviceIds []string, heading string, content string) {

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

/* End File */
