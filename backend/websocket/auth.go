//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Authenticate Connection
//
func (t *Controller) AuthenticateConnection(conn *WebsocketConnection, accessToken string, device_id string) {

	// log connection
	services.Info("Connected Device Id : " + device_id)

	// Store the device id
	conn.muDeviceId.Lock()
	conn.deviceId = device_id
	conn.muDeviceId.Unlock()

	// See if this session is in our db.
	session, err := t.DB.GetByAccessToken(accessToken)

	if err != nil {
		services.Critical("Access Token Not Found - Unable to Authenticate")
		return
	}

	// Get this user is in our db.
	user, err := t.DB.GetUserById(session.UserId)

	if err != nil {
		services.Critical("User Not Found - Unable to Authenticate - UserId : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
		return
	}

	services.Info("Authenticated : " + user.Email)

	// Store the user id from this connection because the auth was successful
	conn.muUserId.Lock()
	conn.userId = user.Id
	conn.muUserId.Unlock()

	// Do the writing.
	go t.DoWsWriting(conn)
}

/* End File */
