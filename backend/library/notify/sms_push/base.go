//
// Date: 2018-10-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-08
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package sms_push

import (
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/sfreiberg/gotwilio"
	"github.com/tidwall/gjson"
)

//
// Push to all sms push channels that the users has opted into.
//
func Push(db models.Datastore, userId uint, uri string, uriRefId uint, data_json string) {

	// Set some helpful times
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Get the status & message.
	title := "Options Cafe Trading"
	status := gjson.Get(data_json, "status").String()
	msg := title + ": The market is now " + status

	// See if we have already sent this notification
	o := models.Notification{}
	db.New().Where("sent_time > ? AND status = ? AND channel = ? AND user_id = ? AND uri = ? AND uri_ref_id = ?", dayStart, "sent", "sms-push", userId, uri, uriRefId).Find(&o)

	// We already sent this notice
	if o.Id > 0 {
		return
	}

	// Switch based on URI
	switch uri {

	case "market-status-open":
		DoMarketStatusChange(db, status, msg)

	case "market-status-closed":
		DoMarketStatusChange(db, status, msg)

	}

	// Store this notification
	ob := models.Notification{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserId:       userId,
		Status:       "sent",
		Channel:      "sms-push",
		Uri:          uri,
		UriRefId:     uint(0),
		Title:        title,
		ShortMessage: "The market is now " + status,
		LongMessage:  "The market is now " + status,
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
func DoMarketStatusChange(db models.Datastore, status string, content string) {

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "postmarket") {
		return
	}

	// TODO: We need to look up in settings if the user wants this notification.

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "sms-push").Find(&nc)

	// Loop through and send message
	for _, row := range nc {
		go DoTwilioSend(row.ChannelId, content)
	}

}

//
// Send message up to twilio
//
func DoTwilioSend(to string, message string) {

	if len(os.Getenv("TWILIO_ACCOUNT_SID")) > 0 {

		// Log
		services.Info("Sending SMD notification: " + to + " - " + message)

		// Setup twilio
		twilio := gotwilio.NewTwilioClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))

		// Send SMS
		twilio.SendSMS(os.Getenv("TWILIO_PHONE"), to, message, "", "")

	}

}

/* End File */
