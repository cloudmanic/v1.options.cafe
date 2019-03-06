//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-12
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package jobs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/pull"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/market"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/worker"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
	nsq "github.com/nsqio/go-nsq"
)

var (
	jobActions        map[string]func(worker.JobRequest) error
	brokerFeedActions map[string]func(models.Datastore, brokers.Api, models.User, models.Broker) error
)

//
// Init
//
func init() {

	// Build standard every day job
	jobActions = map[string]func(worker.JobRequest) error{
		"get-market-status": market.GetMarketStatus,
	}

	// Build out the broker feed action functions
	brokerFeedActions = map[string]func(models.Datastore, brokers.Api, models.User, models.Broker) error{
		"get-orders":              pull.DoGetOrders,
		"get-all-orders":          pull.DoGetAllOrders,
		"get-quotes":              pull.DoGetQuotes,
		"get-balances":            pull.DoGetBalances,
		"get-user-profile":        pull.DoGetUserProfile,
		"get-history":             pull.DoGetHistory,
		"get-positions":           pull.DoGetPositions,
		"do-access-token-refresh": pull.DoAccessTokenRefresh,
		"prime-screener-caches":   screener.PrimeScreenerCachesByUser,
	}

}

//
// Start - Consume messages to make poll a broker
//
func Start(db models.Datastore) {
	// WaitGroup
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()

	config := nsq.NewConfig()

	// New consumer hander
	q, _ := nsq.NewConsumer("oc-job", "oc-worker", config)

	// Conection handler.
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		go HandleRequest(db, string(message.Body))
		return nil
	}))

	// Connect
	err := q.ConnectToNSQD(os.Getenv("NSQD_HOST"))

	if err != nil {
		log.Panic("Could not connect")
	}

	services.Critical("Worker Started...")

	// Wait for messages
	wg.Wait()
}

//
// HandleRequest - Handle broker feed request.
//
func HandleRequest(db models.Datastore, msg string) {
	// Convert JSON to Struct
	job := worker.JobRequest{DB: db}

	if err := json.Unmarshal([]byte(msg), &job); err != nil {
		services.BetterError(err)
		return
	}

	// Is this a normal job
	if _, ok := jobActions[job.Action]; ok {

		// Based on the action sent in call broker API.
		err := jobActions[job.Action](job)

		if err != nil {
			services.BetterError(err)
			return
		}

		return

	}

	// Is this a broker feed request
	if _, ok := brokerFeedActions[job.Action]; ok {

		// Figure out which broker API to use and return it.
		// TODO: we should cache this some how so we are not slamming the DB.
		// And no CPU issues with decrypting the access token. However it is sort
		// of complex because we have the refresh access token call. And access tokens
		// can change at anytime. The cache needs to be busted should an access token
		// change.
		api, user, broker, err := GetBrokerFeedParms(db, job.UserId, job.BrokerId)

		if err != nil {
			services.BetterError(err)
			return
		}

		// Based on the action sent in call broker API.
		err = brokerFeedActions[job.Action](db, api, user, broker)

		if err != nil {
			services.BetterError(err)
			return
		}

		return
	}

	// Return happy
	return
}

//
// Return a broker API object and other data
//
func GetBrokerFeedParms(db models.Datastore, userId uint, brokerId uint) (brokers.Api, models.User, models.Broker, error) {

	var brokerApi brokers.Api

	// Get the user
	user, err := db.GetUserById(userId)

	if err != nil {
		return brokerApi, models.User{}, models.Broker{}, err
	}

	// Get broker
	broker, err := db.GetBrokerById(brokerId)

	if err != nil {
		return brokerApi, models.User{}, models.Broker{}, err
	}

	// Need an access token to continue
	if len(broker.AccessToken) <= 0 {
		return brokerApi, models.User{}, models.Broker{}, errors.New("User Connection (Brokers) No Access Token Found : " + user.Email + " (" + broker.Name + ")")
	}

	// Decrypt the access token
	accessToken, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		return brokerApi, models.User{}, models.Broker{}, err
	}

	// Figure out which broker connection to setup.
	switch broker.Name {

	case "Tradier":
		brokerApi = &tradier.Api{ApiKey: accessToken, DB: db, Sandbox: false}

	case "Tradier Sandbox":
		brokerApi = &tradier.Api{ApiKey: accessToken, DB: db, Sandbox: true}

	default:
		return brokerApi, models.User{}, models.Broker{}, errors.New("Unknown Broker : " + broker.Name + " (" + user.Email + ")")

	}

	// Return happy
	return brokerApi, user, broker, nil
}

/* End File */
