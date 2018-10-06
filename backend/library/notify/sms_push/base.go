//
// Date: 2018-10-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package sms_push

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/sfreiberg/gotwilio"
	"github.com/tidwall/gjson"
)

//
// Push to all sms push channels that the users has opted into.
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

	// Get the status.
	status := gjson.Get(data_json, "status").String()

	// We do not care about the pre-market stuff
	if (status != "open") && (status != "postmarket") {
		return
	}

	// TODO: We need to look up in settings if the user wants this notification.

	// Lets get a list of device ids to send this notification to.
	nc := []models.NotifyChannel{}
	db.New().Where("type = ?", "Sms").Find(&nc)

	// Loop through and send message
	for _, row := range nc {
		go DoTwilioSend(row.ChannelId, "Options Cafe Trading: The market is now "+status)
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
