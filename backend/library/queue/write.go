//
// Date: 2018-11-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package queue

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	nsq "github.com/nsqio/go-nsq"
)

var (
	nsqConn *nsq.Producer
)

//
// Init
//
func init() {

	// NSQ config
	config := nsq.NewConfig()

	// Build producer object
	c, err := nsq.NewProducer(os.Getenv("NSQD_HOST"), config)

	if err != nil {
		services.FatalMsg(err, "PollOrders (NewProducer): NSQ Could not connect.")
	}

	// Set package global
	nsqConn = c
}

//
// Write a message to the queue. We do not
// validate that the message really got into
// the queue because our queue is not safe.
// We assume a message that goes into the queue
// could be lost. In fact they are lost on every deploy.
//
func Write(topic string, msg string) {

	// Send message to message queue
	err := nsqConn.Publish(topic, []byte(msg))

	if err != nil {
		services.FatalMsg(err, "queue.Write: NSQ Could not connect. - "+topic+" : "+msg)
	}

}

/* End File */
