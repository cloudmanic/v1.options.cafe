//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"net/http"

	"app.options.cafe/library/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//
// Check Origin
//
func (t *Controller) CheckOrigin(r *http.Request) bool {

	origin := r.Header.Get("Origin")

	if origin == "file://" {
		return true
	}

	if origin == "http://localhost:4200" {
		return true
	}

	if origin == "http://localhost:7652" {
		return true
	}

	if origin == "https://app.options.cafe" {
		return true
	}

	return false
}

//
// Handle new connections to the app.
//
func (t *Controller) DoWebsocketConnection(c *gin.Context) {

	// setup upgrader
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     t.CheckOrigin,
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		services.Info(err)
		return
	}

	services.InfoMsg("New Websocket Connection - Standard")

	// Close connection when this function ends
	defer conn.Close()

	// Add the connection to our connection array
	r_con := WebsocketConnection{connection: conn, WriteChan: make(chan string, 100)}

	t.Connections[conn] = &r_con

	// Start handling reading messages from the client.
	t.DoWsReading(&r_con)
}

/* End File */
