//
// Date: 2018-11-08
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-08
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"encoding/json"
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Dump the json response for a user's order history
//
// go run main.go --cmd="broker-json-history" --user_id=1 --broker_account_id=2 > orderhistory.json
//
func GetJsonBrokerHistory(db *models.DB, userId int, brokerAccountId int) {

	var brokerApi brokers.Api

	// Get user.
	user, err := db.GetUserById(uint(userId))

	if err != nil {
		panic(err)
	}

	// Get broker account
	brokerAccount, err := db.GetBrokerAccountByIdUserId(uint(brokerAccountId), user.Id)

	if err != nil {
		panic(err)
	}

	// Get broker
	broker, err := db.GetBrokerById(brokerAccount.BrokerId)

	// Decrypt the access token
	accessToken, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		panic(err)
	}

	// Figure out which broker connection to setup.
	switch broker.Name {

	case "Tradier":
		brokerApi = &tradier.Api{ApiKey: accessToken, DB: db, Sandbox: false}

	case "Tradier Sandbox":
		brokerApi = &tradier.Api{ApiKey: accessToken, DB: db, Sandbox: true}

	default:
		panic("Unknown Broker : " + broker.Name + " (" + user.Email + ")")

	}

	// Make API call to get the JSON.
	orders, err := brokerApi.GetAllOrders()

	if err != nil {
		panic(err)
	}

	// Convert to json
	json, err := json.Marshal(orders)

	if err != nil {
		panic(err)
	}

	// Return JSON
	fmt.Println(string(json))
}

/* End File */
