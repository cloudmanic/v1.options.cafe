//
// Date: 2/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package db_feed

import (
	"encoding/json"

	"github.com/cloudmanic/app.options.cafe/backend/controllers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Feed struct {
	DB       *models.DB
	User     models.User
	DataChan chan controllers.SendStruct
}

//
// Start the feed
//
func (t *Feed) Start() {
	services.Info("Starting DB Feed For : " + t.User.Email)

	// Start different DB feeds.
	go t.DoTradeGroupsOpen()
}

//
// Build json to send up websocket.
//
func (t *Feed) WsSend(uri string, data_json string) {

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
		services.BetterError(err)
	}

	t.DataChan <- controllers.SendStruct{UserId: t.User.Id, Body: string(send_json)}
}

/* End File */
