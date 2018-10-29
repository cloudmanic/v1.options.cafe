//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"encoding/json"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/feed"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Start up our user feeds.
//
func (t *Base) StartFeeds() {

	// Get all users (regardless of status)
	users := t.DB.GetAllUsers()

	// Loop through the users
	for i := range users {
		t.DoUserFeed(users[i])
	}

	// Listen for user actions
	go t.ListenForUserActions()

	// Expire Trials
	go t.DoExpireTrails()

	// Do expired Sessions
	go t.ClearExpiredSessions()
}

//
// Clear expired sessions (access tokens)
//
func (t *Base) ClearExpiredSessions() {

	for {

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

		// Sleep for 12 hours
		time.Sleep(time.Second * 60 * 60 * 12)

	}

}

//
// Expire users from Trials
//
func (t *Base) DoExpireTrails() {

	for {

		users := []models.User{}
		t.DB.New().Where("trial_expire <= ? AND status = ? AND stripe_subscription = ?", time.Now(), "Trial", "").Find(&users)

		for _, row := range users {

			row.Status = "Expired"
			t.DB.New().Save(&row)
			services.Info("Free trial has just expired : " + row.Email)
			go services.SlackNotify("#events", "New Options Cafe User Free Trial Expired : "+row.Email)
			go services.SendyUnsubscribe("trial", row.Email)
			go services.SendySubscribe("expired", row.Email, row.FirstName, row.LastName, "", "", "")

		}

		// Sleep for 60 second
		time.Sleep(time.Second * 60)
	}

}

//
// Listen for actions to change the user feed.
//
func (t *Base) ListenForUserActions() {

	for {

		// Wait for action
		msg := <-t.ActionChan

		// Get user
		user, err := t.DB.GetUserById(msg.UserId)

		// Make sure we got the user
		if err != nil {
			services.BetterError(err)
			continue
		}

		// If the action is to restart we restart the feed.
		if msg.Action == "restart" {

			// First we stop the user feed
			t.StopUserFeed(user)

			// Start the feed back up.
			t.DoUserFeed(user)

		}

	}

}

//
// Stop a user feed
//
func (t *Base) StopUserFeed(user models.User) {

	// Make sure there is even a user feed
	if _, ok := t.Users[user.Id]; ok {

		services.Info("Stopping User Connection : " + user.Email)

		// Stop the different broker feeds
		for key := range t.Users[user.Id].BrokerFeed {

			t.Users[user.Id].BrokerFeed[key].MuPolling.Lock()
			t.Users[user.Id].BrokerFeed[key].Polling = false
			t.Users[user.Id].BrokerFeed[key].MuPolling.Unlock()

		}

		// Delete the user map
		delete(t.Users, user.Id)

	}

}

//
// Start one user.
//
func (t *Base) DoUserFeed(user models.User) {

	var brokerApi brokers.Api

	services.Info("Starting User Connection : " + user.Email)

	// This should not happen. But we double check this user is not already started.
	if _, ok := t.Users[user.Id]; ok {
		services.Critical("User Connection Is Already Going : " + user.Email)
		return
	}

	// Verify some default data.
	t.VerifyDefaultWatchList(user)

	// Set the user to the object
	t.Users[user.Id] = &UserFeed{
		Profile:     user,
		WsWriteChan: t.WsWriteChan,
		BrokerFeed:  make(map[uint]*feed.Base),
	}

	// Loop through the different brokers for this user
	for _, row := range t.Users[user.Id].Profile.Brokers {

		// Skip over disabled brokers
		if row.Status == "Disabled" {
			continue
		}

		// Skip over expired brokers
		if row.Status == "Expired" {
			continue
		}

		// Skip over delinquent brokers
		if row.Status == "Delinquent" {
			continue
		}

		// Need an access token to continue
		if len(row.AccessToken) <= 0 {
			services.Critical("User Connection (Brokers) No Access Token Found : " + user.Email + " (" + row.Name + ")")
			continue
		}

		// Decrypt the access token
		decryptAccessToken, err := helpers.Decrypt(row.AccessToken)

		if err != nil {
			services.Error(err, "(DoUserFeed) Unable to decrypt message (#1)")
			continue
		}

		// Figure out which broker connection to setup.
		switch row.Name {

		case "Tradier":
			brokerApi = &tradier.Api{ApiKey: decryptAccessToken, DB: t.DB, Sandbox: false}

		case "Tradier Sandbox":
			brokerApi = &tradier.Api{ApiKey: decryptAccessToken, DB: t.DB, Sandbox: true}

		default:
			services.Critical("Unknown Broker : " + row.Name + " (" + user.Email + ")")
			continue

		}

		// Log magic
		services.Info("Setting up to use " + row.Name + " as the broker for " + user.Email)

		// Set the library we use to fetching data from our broker's API
		t.Users[user.Id].BrokerFeed[row.Id] = &feed.Base{
			DB:          t.DB,
			User:        user,
			BrokerId:    row.Id,
			Api:         brokerApi,
			Polling:     true,
			WsWriteChan: t.WsWriteChan,
		}

		// Start fetching data for this user.
		go t.Users[user.Id].BrokerFeed[row.Id].Start()
	}
}

// ---------------- Helper Functions --------------- //

//
// Build json to send up websocket.
//
func (t *Base) WsSendJsonBuild(uri string, data_json []byte) (string, error) {

	type SendStruct struct {
		Uri  string `json:"uri"`
		Body string `json:"body"`
	}

	// Send Object
	send := SendStruct{
		Uri:  uri,
		Body: string(data_json),
	}
	send_json, err := json.Marshal(send)

	if err != nil {
		services.Error(err, "WsSendJsonBuild() json.Marshal")
		return "", err
	}

	return string(send_json), nil
}

/* End File */
