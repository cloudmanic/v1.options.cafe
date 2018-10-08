//
// Date: 2/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//
// Notes: The idea behind this package is we have a central location
// to send notifications. Send in one message and have it sent to many
// many different places. Some examples might be up websockets, or push
// notifications.
//

package notify

import (
	"encoding/json"
	"flag"

	"github.com/cloudmanic/app.options.cafe/backend/library/notify/sms_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify/web_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
)

var websocketChan chan websocket.SendStruct

//
// Set the web socket channel.
//
func SetWebsocketChannel(ch chan websocket.SendStruct) {
	websocketChan = ch
}

//
// Send notification up websocket.
//
func PushWebsocket(db models.Datastore, userId uint, uri string, uriRefId uint, data_json string) {
	Push(db, userId, []string{"websocket"}, uri, uriRefId, data_json)
}

//
// A general way to push to all channels
//
func Push(db models.Datastore, userId uint, channels []string, uri string, uriRefId uint, data_json string) {

	// Do nothing if we are testing. TODO: build testing for this.
	if flag.Lookup("test.v") != nil {
		return
	}

	type SendStruct struct {
		Uri  string `json:"uri"`
		Body string `json:"body"`
	}

	// Send Object
	send := SendStruct{
		Uri:  uri,
		Body: data_json,
	}

	send_json, err := json.Marshal(send)

	if err != nil {
		services.BetterError(err)
	}

	// Loop through and send to all the different channels.
	for _, row := range channels {

		switch row {

		case "websocket":
			websocketChan <- websocket.SendStruct{UserId: userId, Body: string(send_json)}

		case "web-push":
			go web_push.Push(db, userId, uri, uriRefId, data_json)

		case "sms-push":
			go sms_push.Push(db, userId, uri, uriRefId, data_json)

		}

	}

}

/* End File */
