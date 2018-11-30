//
// Date: 2018-11-10
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-25
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package user

import (
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Clear expired sessions (access tokens)
//
func ClearExpiredSessions(db *models.DB) {

	// Find the Centcom app
	centcomApp := models.Application{}
	db.New().Where("name <= ?", "Centcom").Find(&centcomApp)

	if centcomApp.Id > 0 {

		// Delete expired centcom sessions
		db.New().Where("last_activity <= ? AND application_id = ?", time.Now().AddDate(0, 0, -1), centcomApp.Id).Delete(&models.Session{})

		// Just cleared Centcom sessions
		services.Info("Centcom sessions cleared.")

	}

	// Find the Personal app
	personalApp := models.Application{}
	db.New().Where("name <= ?", "Personal Access Token").Find(&personalApp)

	if personalApp.Id > 0 {

		// Clear all sessions that have not had activity in the last 14 days (2 weeks)
		db.New().Where("last_activity <= ? AND application_id != ?", time.Now().AddDate(0, 0, -14), personalApp.Id).Delete(&models.Session{})

		// Log cleared sessions.
		services.Info("All expired sessions cleared.")

	}

}

//
// Expire users from Trials
//
func ExpireTrails(db *models.DB) {

	users := []models.User{}
	db.New().Debug().Where("trial_expire <= ? AND status = ? AND stripe_subscription = ?", time.Now(), "Trial", "").Find(&users)

	for _, row := range users {

		services.Info(strconv.Itoa(int(row.Id)))

		if row.Id <= 0 {
			continue
		}

		row.Status = "Expired"
		db.New().Save(&row)
		services.Info("Free trial has just expired : " + row.Email)
		go services.SlackNotify("#events", "New Options Cafe User Free Trial Expired : "+row.Email)
		go services.SendyUnsubscribe("trial", row.Email)
		go services.SendySubscribe("expired", row.Email, row.FirstName, row.LastName, "", "", "", "No")

	}

	services.Info("All expire trails set to expired.")
}

/* End File */
