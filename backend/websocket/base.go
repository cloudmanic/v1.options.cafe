package websocket

import (
	"sync"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gorilla/websocket"
)

const writeWait = 10 * time.Second

type Controller struct {
	DB          models.Datastore
	WsWriteChan chan SendStruct
	Connections map[*websocket.Conn]*WebsocketConnection
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
	Body   string `json:"body"`
	UserId uint   `json:"user_id"`
}

//
// Create a new instance of a websocket controller.
//
func NewController(db models.Datastore, WsWriteChan chan SendStruct) *Controller {

	return &Controller{
		DB:          db,
		WsWriteChan: WsWriteChan,
		Connections: make(map[*websocket.Conn]*WebsocketConnection),
	}
}

/* End File */
