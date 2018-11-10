//
// Date: 2018-11-10
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-10
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"encoding/json"
	"fmt"
	"os"

	nsq "github.com/bitly/go-nsq"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

var nsqConn *nsq.Producer

type SendStruct struct {
	Uri    string `json:"uri"`
	UserId uint   `json:"user_id"`
	Body   string `json:"body"`
}

//
// Init
//
func init() {

	// NSQ config
	config := nsq.NewConfig()

	// Build producer object
	c, err := nsq.NewProducer(os.Getenv("NSQD_HOST"), config)

	if err != nil {
		services.FatalMsg(err, "poll.init (NewProducer): NSQ Could not connect.")
	}

	// Set package global
	nsqConn = c

}

//
// Send data up websocket.
//
func WriteWebsocket(user models.User, send_type string, sendObject interface{}) error {

	// Convert to a json string.
	dataJson, err := json.Marshal(sendObject)

	if err != nil {
		return fmt.Errorf("pull.WriteWebsocket() json.Marshal : ", err)
	}

	// Send data up websocket.
	sendJson, err := GetSendJson(user, send_type, string(dataJson))

	if err != nil {
		return fmt.Errorf("pull.WriteWebsocket() GetSendJson Send Object : ", err)
	}

	// Send message to message queue
	err = nsqConn.Publish("oc-websocket-write", []byte(sendJson))

	if err != nil {
		services.FatalMsg(err, "WriteWebsocket: NSQ Could not connect.")
	}

	// Return happy.
	return nil
}

//
// Return a json object ready to be sent up the websocket
//
func GetSendJson(user models.User, uri string, data_json string) (string, error) {

	// Send Object
	send := SendStruct{
		Uri:    uri,
		UserId: user.Id,
		Body:   data_json,
	}

	send_json, err := json.Marshal(send)

	if err != nil {
		return "", err
	}

	return string(send_json), nil
}

/* End File */
