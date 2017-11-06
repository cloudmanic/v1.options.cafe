package controllers

import (
	"fmt"

	"app.options.cafe/backend/library/services"
	"github.com/tidwall/gjson"
)

//
// Process a read request that was sent in from the client
//
func (t *Controller) ProcessRead(conn *WebsocketConnection, message string, data map[string]interface{}) {

	switch data["type"] {

	// Ping to make sure we are alive.
	case "ping":
		conn.WriteChan <- "{\"type\":\"pong\"}"
		break

	// // Refresh all cached data.
	// case "refresh-all-data":
	// 	t.WsReadChan <- SendStruct{UserId: conn.userId, Message: "FromCache:refresh"}
	// 	break

	// // Refresh watchlists
	// case "refresh-watchlists":
	// 	t.WsReadChan <- SendStruct{UserId: conn.userId, Message: "Watchlists:refresh"}
	// 	break

	// The user authenticates.
	case "set-access-token":
		device_id := gjson.Get(message, "data.device_id").String()
		access_token := gjson.Get(message, "data.access_token").String()
		t.AuthenticateConnection(conn, access_token, device_id)
		break

	// Default we send over to the user feed.
	default:
		t.WsReadChan <- ReceivedStruct{UserId: conn.userId, Message: message, Connection: conn}
		break

	}

}

//
// Authenticate Connection
//
func (t *Controller) AuthenticateConnection(conn *WebsocketConnection, accessToken string, device_id string) {

	// log connection
	services.Log("Connected Device Id : " + device_id)

	// Store the device id
	conn.muDeviceId.Lock()
	conn.deviceId = device_id
	conn.muDeviceId.Unlock()

	// See if this session is in our db.
	session, err := t.DB.GetByAccessToken(accessToken)

	if err != nil {
		services.MajorLog("Access Token Not Found - Unable to Authenticate")
		return
	}

	// Get this user is in our db.
	user, err := t.DB.GetUserById(session.UserId)

	if err != nil {
		services.MajorLog("User Not Found - Unable to Authenticate - UserId : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
		return
	}

	services.Log("Authenticated : " + user.Email)

	// Store the user id from this connection because the auth was successful
	conn.muUserId.Lock()
	conn.userId = user.Id
	conn.muUserId.Unlock()

	// Do the writing.
	go t.DoWsWriting(conn)

	// Send cached data so they do not have to wait for polling.
	t.WsReadChan <- ReceivedStruct{UserId: conn.userId, Message: "refresh-all-data", Connection: conn}

}

/* End File */
