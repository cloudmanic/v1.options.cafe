//
// Date: 2018-11-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package queue

import (
	"errors"
	"os"

	"app.options.cafe/library/services"
	nsq "github.com/nsqio/go-nsq"
)

var (
	nsqConn *nsq.Producer
)

//
// Just start the connect to the queue.
//
func Start() {
	// NSQ config
	config := nsq.NewConfig()

	// Build producer object
	c, err := nsq.NewProducer(os.Getenv("NSQD_HOST"), config)

	if err != nil {
		services.Fatal(errors.New(err.Error() + "PollOrders (NewProducer): NSQ Could not connect."))
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
	// We have not started the message queue. Maybe in testing.
	if nsqConn == nil {
		return
	}

	// Send message to message queue
	err := nsqConn.Publish(topic, []byte(msg))

	if err != nil {
		services.Fatal(errors.New(err.Error() + "queue.Write: NSQ Could not connect. - " + topic + " : " + msg))
	}

}

/* End File */
