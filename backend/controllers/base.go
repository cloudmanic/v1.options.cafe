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

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const defaultMysqlLimit = 100
const httpNoRecordFound = "No Record Found."
const httpGenericErrMsg = "Please contact support at help@options.cafe."

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
// RespondJSON makes the response with payload as json format.
// This is used when we want the json back (used in websockets).
// If you do not need the json back just use c.JSON()
//
func (t *Controller) RespondJSON(c *gin.Context, status int, payload interface{}) string {

	response, err := json.Marshal(payload)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return ""
	}

	// Return json.
	c.JSON(200, payload)

	// We return the raw JSON
	return string(response)
}

//
// Return error.
//
func (t *Controller) RespondError(c *gin.Context, err error, msg string) bool {

	if err != nil {
		services.Warning(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return true
	}

	// No error.
	return false
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

/* End File */
