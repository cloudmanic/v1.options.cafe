//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"app.options.cafe/backend/models"
	"github.com/gorilla/websocket"
)

type Controller struct {
	DB                models.Datastore
	WsReadChan        chan SendStruct
	WsWriteChan       chan SendStruct
	WsWriteQuoteChan  chan SendStruct
	Connections       map[*websocket.Conn]*WebsocketConnection
	QuotesConnections map[*websocket.Conn]*WebsocketConnection
}

/* End File */
