//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"sync"

	"app.options.cafe/backend/models"
	"github.com/gorilla/websocket"
)

type Controller struct {
	DB                models.Datastore
	WsReadChan        chan ReceivedStruct
	WsWriteChan       chan SendStruct
	WsWriteQuoteChan  chan SendStruct
	Connections       map[*websocket.Conn]*WebsocketConnection
	QuotesConnections map[*websocket.Conn]*WebsocketConnection
}

type WebsocketConnection struct {
	writeChan  chan string
	connection *websocket.Conn

	muUserId sync.Mutex
	userId   uint

	muDeviceId sync.Mutex
	deviceId   string
}

type SendStruct struct {
	UserId  uint
	Message string
}

type ReceivedStruct struct {
	UserId     uint
	Message    string
	Connection *WebsocketConnection
}

/* End File */
