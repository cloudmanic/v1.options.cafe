//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"encoding/json"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gorilla/websocket"
)

//
// Start a writer for the websocket connection.
//
func (t *Controller) DoWsWriting(conn *WebsocketConnection) {

	conn.connection.SetWriteDeadline(time.Now().Add(writeWait))

	for {
		message := <-conn.WriteChan
		conn.connection.WriteMessage(websocket.TextMessage, []byte(message))
		conn.connection.SetWriteDeadline(time.Now().Add(writeWait))
	}
}

//
// Send a message to all connected clients.
//
func (t *Controller) WsDispatchToAll(send string) {
	for i := range t.Connections {
		t.Connections[i].WriteChan <- send
	}
}

//
// Listen for data from our broker feeds.
// Take the data and then pass it up the websockets.
//
func (t *Controller) DoWsDispatch() {

	for {

		send := <-t.WsWriteChan

		for i := range t.Connections {

			// We only care about the user we passed in.
			if t.Connections[i].userId == send.UserId {
				t.Connections[i].WriteChan <- send.Body
			}

		}

	}

}

//
// Build json to send up websocket.
//
func (t *Controller) WsSendJsonBuild(uri string, data_json string) (string, error) {

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
		return "", err
	}

	return string(send_json), nil
}

/* End File */
