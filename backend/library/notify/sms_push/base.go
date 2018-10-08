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

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/sfreiberg/gotwilio"
)

//
// Push to all sms push channels that the users has opted into.
//
func Push(to string, message string) {
	DoTwilioSend(to, message)
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
