//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package polling

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

var (

	// Here we set the different type of polls and how often we poll them
	polls []Poll = []Poll{
		{Name: "get-market-status", Sleep: 2, Delay: 0, Type: "simple"},
		{Name: "get-quotes", Sleep: 1, Delay: 0, Type: "all-users"},
		{Name: "get-orders", Sleep: 3, Delay: 30, Type: "all-users"},        // Give get-all-orders enough time to get started
		{Name: "get-all-orders", Sleep: 86400, Delay: 0, Type: "all-users"}, // 24 hours
		{Name: "get-balances", Sleep: 5, Delay: 0, Type: "all-users"},
		{Name: "get-user-profile", Sleep: 20, Delay: 0, Type: "all-users"},
		{Name: "get-history", Sleep: 43200, Delay: 0, Type: "all-users"},   // 12 hours
		{Name: "get-positions", Sleep: 3600, Delay: 10, Type: "all-users"}, // 1 hour, 10 seconds delay (we want all orders to complete first)
		{Name: "do-access-token-refresh", Sleep: 60, Delay: 0, Type: "all-users"},
		{Name: "prime-screener-caches", Sleep: 240, Delay: 0, Type: "all-users"}, // 4 mins
	}
)

type Poll struct {
	Name  string
	Type  string
	Delay time.Duration // seconds
	Sleep time.Duration // seconds
}

//
// Start broker polling. We really should only have
// one machine that does the broker polling.
// Via the ENV files we determine if this machine should
// do the broker polling or not.
//
func Start(db models.Datastore) {

	// Small delay to make sure all services are up.
	time.Sleep(time.Second * 3)

	// Start different types of polls.
	for _, row := range polls {
		go StartPoll(db, row)
	}

	// When we call this from the CLI via "-cmd=broker-feed-poller". Only used in production.
	// GoExit will exit after all go routines exit
	if os.Getenv("APP_ENV") != "local" {
		runtime.Goexit()
	}

	// Log
	services.InfoMsg("Broker Feed Poller Started...")
}

//
// Start a poll
//
func StartPoll(db models.Datastore, poll Poll) {

	services.InfoMsg("Starting polling for " + poll.Name + ".")

	// Delay before starting
	if poll.Delay > 0 {
		time.Sleep(time.Second * poll.Delay)
	}

	// Start polling
	for {

		switch poll.Type {

		// Send action to all users
		case "simple":
			SendSimpleAction(poll.Name)

		// Just send one request
		case "all-users":
			SendActionToAllUsers(db, poll.Name)

		}

		// Sleep for X seconds
		time.Sleep(time.Second * poll.Sleep)
	}

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
// Send a simple message
//
func SendSimpleAction(action string) {
	// Send message to message queue
	queue.Write("oc-job", `{"action":"`+action+`"}`)
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

			// Send message to message queue
			queue.Write("oc-job", `{"action":"`+action+`","user_id":`+strconv.Itoa(int(row.Id))+`,"broker_id":`+strconv.Itoa(int(row2.Id))+`}`)

		}

	}

}

/* End File */
