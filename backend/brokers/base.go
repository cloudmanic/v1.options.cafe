//
// Date: 2018-11-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package brokers

import (
	"errors"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Sometimes we have generic needs. Such as historical stock quotes.
// We use this to pick the more optimal broker to support this.
// This function will return the optimal Tradier broker connection.
//
func GetPrimaryTradierConnection(db models.Datastore, userId uint) (tradier.Api, error) {

	// Get the full user
	user, err := db.GetUserById(userId)

	if err != nil {
		return tradier.Api{}, err
	}

	// Get a list of brokers this user has.
	brokers := []models.Broker{}

	// Run the query to get brokers. For now we always get Tradier
	db.New().Where("user_id = ?", user.Id).Find(&brokers)

	// If we have no brokers return error : TODO: change this to default to our Tradier ADMIN Key
	if len(brokers) <= 0 {
		return tradier.Api{}, err
	}

	// Find our default broker
	broker := models.Broker{}

	for _, row := range brokers {

		if (row.Name == "Tradier") || (row.Name == "Tradier Sandbox") {
			broker = row
			break
		}

	}

	// We did not find a good broker : TODO: change this to default to our Tradier ADMIN Key
	if broker.Id <= 0 {
		return tradier.Api{}, errors.New("No good broker found.")
	}

	// Decrypt the access token
	accessToken, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		panic(err)
	}

	// Figure out which broker connection to setup.
	var brokerApi tradier.Api

	switch broker.Name {

	case "Tradier":
		brokerApi = tradier.Api{ApiKey: accessToken, DB: db, Sandbox: false}

	case "Tradier Sandbox":
		brokerApi = tradier.Api{ApiKey: accessToken, DB: db, Sandbox: true}

	default:
		return tradier.Api{}, errors.New("No good broker found.")

	}

	// Return Happy
	return brokerApi, nil

}

/* End File */
