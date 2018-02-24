package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const writeWait = 5 * time.Second

type Controller struct {
	DB                models.Datastore
	WsWriteChan       chan SendStruct
	WsWriteQuoteChan  chan SendStruct
	Connections       map[*websocket.Conn]*WebsocketConnection
	QuotesConnections map[*websocket.Conn]*WebsocketConnection
}

type WebsocketConnection struct {
	WriteChan  chan string
	connection *websocket.Conn

	muUserId sync.Mutex
	userId   uint

	muDeviceId sync.Mutex
	deviceId   string
}

type SendStruct struct {
	Body   string
	UserId uint
}

//
// Create a new instance of a websocket controller.
//
func NewController(db models.Datastore, WsWriteChan chan SendStruct, WsWriteQuoteChan chan SendStruct) *Controller {

	return &Controller{
		DB:                db,
		WsWriteChan:       WsWriteChan,
		WsWriteQuoteChan:  WsWriteQuoteChan,
		Connections:       make(map[*websocket.Conn]*WebsocketConnection),
		QuotesConnections: make(map[*websocket.Conn]*WebsocketConnection),
	}
}

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
		services.BetterError(err)
		return
	}

	services.Info("New Websocket Connection - Standard")

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
func (t *Controller) DoQuoteWebsocketConnection(c *gin.Context) {

	// setup upgrader
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     t.CheckOrigin,
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		services.BetterError(err)
		return
	}

	services.Info("New Websocket Connection - Quote")

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

//
// Listen for data from our broker feeds.
// Take the data and then pass it up the websockets.
//
func (t *Controller) DoWsDispatch() {

	for {

		select {

		// Core channel
		case send := <-t.WsWriteChan:

			for i := range t.Connections {

				// We only care about the user we passed in.
				if t.Connections[i].userId == send.UserId {

					select {

					case t.Connections[i].WriteChan <- send.Body:

					default:
						services.Critical("Channel full. Discarding value (Core channel)")

					}

				}

			}

		// Quotes channel
		case send := <-t.WsWriteQuoteChan:

			for i := range t.QuotesConnections {

				// We only care about the user we passed in.
				if t.QuotesConnections[i].userId == send.UserId {

					select {

					case t.QuotesConnections[i].WriteChan <- send.Body:

					default:
						services.Critical("Channel full. Discarding value (Quotes channel)")

					}

				}

			}

		}

	}

}

/* End File */
