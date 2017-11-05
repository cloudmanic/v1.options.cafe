//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"encoding/json"

	"app.options.cafe/backend/brokers"
	"app.options.cafe/backend/brokers/feed"
	"app.options.cafe/backend/brokers/tradier"
	"app.options.cafe/backend/controllers"
	"app.options.cafe/backend/library/helpers"
	"app.options.cafe/backend/library/services"
	"app.options.cafe/backend/models"
)

type UserFeed struct {
	Profile    models.User
	DataChan   chan controllers.SendStruct
	BrokerFeed map[uint]*feed.Base
}

//
// Start up our user feeds.
//
func (t *Base) StartFeeds() {

	// Get all active users
	users := t.DB.GetAllActiveUsers()

	// Loop through the users
	for i := range users {
		t.DoUserFeed(users[i])
	}

	// Listen of income Feed Requests.
	go t.DoFeedRequestListen()

}

//
// Start one user.
//
func (t *Base) DoUserFeed(user models.User) {

	var brokerApi brokers.Api

	services.Log("Starting User Connection : " + user.Email)

	// This should not happen. But we double check this user is not already started.
	if _, ok := t.Users[user.Id]; ok {
		services.MajorLog("User Connection Is Already Going : " + user.Email)
		return
	}

	// Verify some default data.
	t.VerifyDefaultWatchList(user)

	// Set the user to the object
	t.Users[user.Id] = &UserFeed{
		Profile:    user,
		DataChan:   t.DataChan,
		BrokerFeed: make(map[uint]*feed.Base),
	}

	// Loop through the different brokers for this user
	for _, row := range t.Users[user.Id].Profile.Brokers {

		// Need an access token to continue
		if len(row.AccessToken) <= 0 {
			services.MajorLog("User Connection (Brokers) No Access Token Found : " + user.Email + " (" + row.Name + ")")
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
			brokerApi = &tradier.Api{ApiKey: decryptAccessToken, DB: t.DB}

		default:
			services.MajorLog("Unknown Broker : " + row.Name + " (" + user.Email + ")")
			continue

		}

		// Log magic
		services.Log("Setting up to use " + row.Name + " as the broker for " + user.Email)

		// Set the library we use to fetching data from our broker's API
		t.Users[user.Id].BrokerFeed[row.Id] = &feed.Base{
			DB:        t.DB,
			User:      user,
			Api:       brokerApi,
			DataChan:  t.DataChan,
			QuoteChan: t.QuoteChan,
		}

		// Start fetching data for this user.
		go t.Users[user.Id].BrokerFeed[row.Id].Start()
	}
}

//
// Refresh all data.
//
func (t *Base) RefreshAllData(user *UserFeed) {

	// Loop through each broker and refresh the data.
	for _, row := range t.Users[user.Profile.Id].BrokerFeed {
		row.RefreshFromCached()
	}

	// Send watchlist
	t.WsSendWatchlists(t.Users[user.Profile.Id])
}

// ---------------- Helper Functions --------------- //

//
// Build json to send up websocket.
//
func (t *Base) WsSendJsonBuild(send_type string, data_json []byte) (string, error) {

	type SendStruct struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}

	// Send Object
	send := SendStruct{
		Type: send_type,
		Data: string(data_json),
	}
	send_json, err := json.Marshal(send)

	if err != nil {
		services.Error(err, "WsSendJsonBuild() json.Marshal")
		return "", err
	}

	return string(send_json), nil
}

/* End File */
