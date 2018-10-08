//
// Date: 2018-10-08
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-08
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket_push

import (
	"encoding/json"
	"flag"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
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
// A general way to push to all channels
//
func Push(userId uint, uri string, data_json string) {

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

	// Send up websocket
	websocketChan <- websocket.SendStruct{UserId: userId, Body: string(send_json)}
}

/* End File */
