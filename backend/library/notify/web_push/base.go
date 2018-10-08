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

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	onesignal "github.com/tbalthazar/onesignal-go"
)

//
// Push to all web push channels that the users has opted into.
//
func Push(deviceIds []string, heading string, content string) {
	DoOneSignalWebPushSend(deviceIds, heading, content)
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
