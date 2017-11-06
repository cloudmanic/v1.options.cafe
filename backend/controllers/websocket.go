package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"app.options.cafe/backend/library/services"
	"github.com/gorilla/websocket"
)

const writeWait = 5 * time.Second

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
func (t *Controller) DoWebsocketConnection(w http.ResponseWriter, r *http.Request) {

	// setup upgrader
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     t.CheckOrigin,
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		services.Error(err, "Upable to upgrade Websocket")
		return
	}

	services.Log("New Websocket Connection - Standard")

	// Close connection when this function ends
	defer conn.Close()

	// Add the connection to our connection array
	r_con := WebsocketConnection{connection: conn, WriteChan: make(chan string, 100)}

	t.Connections[conn] = &r_con

	// Start handling reading messages from the client.
	t.DoWsReading(&r_con)
}

//
// Handle new quote connections to the app.
//
func (t *Controller) DoQuoteWebsocketConnection(w http.ResponseWriter, r *http.Request) {

	// setup upgrader
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     t.CheckOrigin,
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		services.Error(err, "(DoQuoteWebsocketConnection) Unable to upgrade Quote Websocket")
		return
	}

	services.Log("New Websocket Connection - Quote")

	// Close connection when this function ends
	defer conn.Close()

	// Add the connection to our connection array
	r_con := WebsocketConnection{connection: conn, WriteChan: make(chan string, 1000)}

	t.QuotesConnections[conn] = &r_con

	// Do reading
	t.DoWsReading(&r_con)
}

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
				services.Log("Client Disconnected (" + conn.deviceId + ") ...")
			}

			_, ok2 := t.QuotesConnections[conn.connection]

			if ok2 {
				delete(t.QuotesConnections, conn.connection)
				services.Log("Client Quote Disconnected (" + conn.deviceId + ") ...")
			}

			break
		}

		// this should come after the mt test.
		if err != nil {
			services.Error(err, "Error in DoWsReading")
			break
		}

		// Json decode message.
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			services.Error(err, "(DoWsReading) Unable to decode json")
			break
		}

		// Switch on the type of requests.
		t.ProcessRead(conn, string(message), data)

	}
}

/* End File */
