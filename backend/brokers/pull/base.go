//
// Date: 2019-04-13
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"strconv"
	"strings"
	"time"

	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// CheckErrorForDisableBroker - If an API call fails we check if we should
// disable the broker and then tell the user to reconnect.
//
func CheckErrorForDisableBroker(loc string, err error, db models.Datastore, user models.User, broker models.Broker) {

	services.InfoMsg(loc + ": Error with User: " + strconv.Itoa(int(user.Id)) + " : " + err.Error())

	// If the error has this string we disable the broker.
	if strings.Contains(err.Error(), "Invalid Access Token") {
		// Disable broker
		broker.Status = "Disabled"
		broker.AccessToken = ""
		broker.RefreshToken = ""
		db.New().Save(&broker)

		// Log notice for dashboard
		// Store this notification
		const title = "Lost connection to your broker"

		ob := models.Notification{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserId:      user.Id,
			Status:      "pending",
			Channel:     "in-app",
			Uri:         "dashboard-notice",
			Title:       title,
			LongMessage: `It seems we have lost your connection to Tradier. Please reconnect by visiting our <a href="/settings/brokers" target="_self">broker section</a>.`,
			SentTime:    time.Now(),
			Expires:     time.Now().AddDate(0, 10, 0), // 10 months
		}

		// Make sure there is not already a notice like this in the DB.
		n := models.Notification{}
		db.New().Where("status = ? AND user_id = ? AND title = ? AND channel = ?", "pending", user.Id, title, "in-app").Find(&n)

		// Store in DB
		if n.Id == 0 {
			db.New().Save(&ob)
			services.InfoMsg("Disabled broker account " + strconv.Itoa(int(broker.Id)) + " for user " + strconv.Itoa(int(user.Id)) + " due to invalid access token")
		}

		// TODO(spicer): Send user an email.
	}
}
