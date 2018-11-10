//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-09
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package polling

import (
	"os"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	nsq "github.com/nsqio/go-nsq"
)

var nsqConn *nsq.Producer

//
// Init
//
func init() {

	// NSQ config
	config := nsq.NewConfig()

	// Build producer object
	c, err := nsq.NewProducer(os.Getenv("NSQ_HOST"), config)

	if err != nil {
		services.FatalMsg(err, "PollOrders (NewProducer): NSQ Could not connect.")
	}

	// Set package global
	nsqConn = c
}

//
// Start broker polling
//
func Start(db models.Datastore) {
	go PollOrders(db)
}

//
// Return a list of active users.
//
func GetActiveUserList(db models.Datastore) []models.User {

	// Get all users
	users := db.GetAllActiveOrTrialUsers()

	// TODO put some caching in here so we do not slam the DB.

	// Return the user list
	return users
}

//
// Send message out to all users.
//
func SendActionToAllUsers(db models.Datastore, action string) {

	// Get a list of users we are polling for
	users := GetActiveUserList(db)

	// Lopp through the users we just got.
	for _, row := range users {

		// Loop through the different brokers
		for _, row2 := range row.Brokers {

			// Make sure the broker is active
			if row2.Status != "Active" {
				continue
			}

			// Build raw json
			rawJson := `{"action":"` + action + `","user_id":` + strconv.Itoa(int(row.Id)) + `,"broker_id":` + strconv.Itoa(int(row2.Id)) + `}`

			// Send message to message queue
			err := nsqConn.Publish("oc-broker-feed-request", []byte(rawJson))

			if err != nil {
				services.FatalMsg(err, "SendActionToAllUsers: NSQ Could not connect. - "+action)
			}

		}

	}

}

/* End File */
