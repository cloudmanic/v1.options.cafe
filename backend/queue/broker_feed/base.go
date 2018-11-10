//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-10
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package broker_feed

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
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	nsq "github.com/nsqio/go-nsq"
)

var (
	brokerFeedActions map[string]func(models.Datastore, brokers.Api, models.User, models.Broker) error
)

type ActionRequest struct {
	Action   string `json:"action"`
	UserId   uint   `json:"user_id"`
	BrokerId uint   `json:"broker_id"`
}

//
// Init
//
func init() {

	// Build out the action functions
	brokerFeedActions = map[string]func(models.Datastore, brokers.Api, models.User, models.Broker) error{
		"get-orders": pull.DoGetOrders,
		"get-quotes": pull.DoGetQuotes,
	}

}

//
// Consume messages to make poll a broker
//
func Start(db models.Datastore) {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()

	config := nsq.NewConfig()

	// New consumer hander
	q, _ := nsq.NewConsumer("oc-broker-feed-request", "oc-broker-feed-worker", config)

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

	services.Critical("Broker Feed Worker Started...")

	// Wait for messages
	wg.Wait()

}

//
// Handle broker feed request.
//
func HandleRequest(db models.Datastore, msg string) {

	// Convert JSON to Struct
	ac := ActionRequest{}

	if err := json.Unmarshal([]byte(msg), &ac); err != nil {
		services.BetterError(err)
		return
	}

	// Figure out which broker API to use and return it.
	// TODO: we should cache this some how so we are not slamming the DB.
	// And no CPU issues with decrypting the access token. However it is sort
	// of complex because we have the refresh access token call. And access tokens
	// can change at anytime. The cache needs to be busted should an access token
	// change.
	api, user, broker, err := GetBrokerFeedParms(db, ac.UserId, ac.BrokerId)

	if err != nil {
		services.BetterError(err)
		return
	}

	// Make sure the action we went in is known
	if _, ok := brokerFeedActions[ac.Action]; !ok {
		services.BetterError(errors.New("Unknown broker feed action."))
		return
	}

	// Based on the action sent in call broker API.
	err = brokerFeedActions[ac.Action](db, api, user, broker)

	if err != nil {
		services.BetterError(err)
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
