//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"fmt"

	"app.options.cafe/library/services"
)

//
// Authenticate Connection
//
func (t *Controller) AuthenticateConnection(conn *WebsocketConnection, accessToken string, device_id string) {

	// log connection
	services.InfoMsg("Connected Device Id : " + device_id)

	// Store the device id
	conn.muDeviceId.Lock()
	conn.deviceId = device_id
	conn.muDeviceId.Unlock()

	// See if this session is in our db.
	session, err := t.DB.GetByAccessToken(accessToken)

	if err != nil {
		services.InfoMsg("Access Token Not Found - Unable to Authenticate")
		return
	}

	// Get this user is in our db.
	user, err := t.DB.GetUserById(session.UserId)

	if err != nil {
		services.InfoMsg("User Not Found - Unable to Authenticate - UserId : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
		return
	}

	services.InfoMsg("Authenticated : " + user.Email)

	// Store the user id from this connection because the auth was successful
	conn.muUserId.Lock()
	conn.userId = user.Id
	conn.muUserId.Unlock()

	// Do the writing.
	go t.DoWsWriting(conn)
}

/* End File */
