//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
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

type ReceivedStruct struct {
	Body       string
	UserId     uint
	Connection *WebsocketConnection
}

//
// RespondJSON makes the response with payload as json format
//
func (t *Controller) RespondJSON(w http.ResponseWriter, status int, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Return json.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

//
// RespondError makes the error response with payload as json format
//
func (t *Controller) RespondError(w http.ResponseWriter, code int, message string) {
	t.RespondJSON(w, code, map[string]string{"error": message})
}

/* End File */
