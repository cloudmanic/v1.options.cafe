//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"encoding/json"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/tidwall/gjson"
)

//
// Process a read request that was sent in from the client
//
func (t *Controller) ProcessRead(conn *WebsocketConnection, message string, data map[string]interface{}) {

	switch data["uri"] {

	// Ping to make sure we are alive.
	case "ping":
		conn.WriteChan <- "{\"uri\":\"pong\"}"
		break

	// The user authenticates.
	case "set-access-token":
		device_id := gjson.Get(message, "body.device_id").String()
		access_token := gjson.Get(message, "body.access_token").String()
		t.AuthenticateConnection(conn, access_token, device_id)
		break

	// Default we send over to the user feed.
	default:
		services.Info("Unknown message coming from websocket - " + message)
		//t.WsReadChan <- ReceivedStruct{UserId: conn.userId, Body: message, Connection: conn}
		break

	}

}

//
// Start a reader for this websocket connection.
//
func (t *Controller) DoWsReading(conn *WebsocketConnection) {

	for {

		// Block waiting for a message to arrive
		mt, message, err := conn.connection.ReadMessage()

		// Connection closed.
		if mt < 0 {

			_, ok := t.Connections[conn.connection]

			if ok {
				delete(t.Connections, conn.connection)
				services.Info("Client Disconnected (" + conn.deviceId + ") ...")
			}

			_, ok2 := t.QuotesConnections[conn.connection]

			if ok2 {
				delete(t.QuotesConnections, conn.connection)
				services.Info("Client Quote Disconnected (" + conn.deviceId + ") ...")
			}

			break
		}

		// this should come after the mt test.
		if err != nil {
			services.BetterError(err)
			break
		}

		// Json decode message.
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			services.BetterError(err)
			break
		}

		// Switch on the type of requests.
		t.ProcessRead(conn, string(message), data)

	}
}

/* End File */
