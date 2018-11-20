//
// Date: 2018-11-10
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-19
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package user

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Base struct {
	DB models.Datastore
}

//
// Clear expired sessions (access tokens)
//
func (t *Base) ClearExpiredSessions() {

	// Find the Centcom app
	centcomApp := models.Application{}
	t.DB.New().Where("name <= ?", "Centcom").Find(&centcomApp)

	if centcomApp.Id > 0 {

		// Delete expired centcom sessions
		t.DB.New().Where("last_activity <= ? AND application_id = ?", time.Now().AddDate(0, 0, -1), centcomApp.Id).Delete(&models.Session{})

		// Just cleared Centcom sessions
		services.Info("Centcom sessions cleared.")

	}

	// Find the Personal app
	personalApp := models.Application{}
	t.DB.New().Where("name <= ?", "Personal Access Token").Find(&personalApp)

	if personalApp.Id > 0 {

		// Clear all sessions that have not had activity in the last 14 days (2 weeks)
		t.DB.New().Where("last_activity <= ? AND application_id != ?", time.Now().AddDate(0, 0, -14), personalApp.Id).Delete(&models.Session{})

		// Log cleared sessions.
		services.Info("All expired sessions cleared.")

	}

}

//
// Expire users from Trials
//
func (t *Base) ExpireTrails() {

	users := []models.User{}
	t.DB.New().Where("trial_expire <= ? AND status = ? AND stripe_subscription = ?", time.Now(), "Trial", "").Find(&users)

	for _, row := range users {

		if row.Id <= 0 {
			continue
		}

		row.Status = "Expired"
		t.DB.New().Save(&row)
		services.Info("Free trial has just expired : " + row.Email)
		go services.SlackNotify("#events", "New Options Cafe User Free Trial Expired : "+row.Email)
		go services.SendyUnsubscribe("trial", row.Email)
		go services.SendySubscribe("expired", row.Email, row.FirstName, row.LastName, "", "", "", "No")

	}

	services.Info("All expire trails set to expired.")
}

/* End File */
