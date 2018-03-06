//
// Date: 3/1/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Create a new application.
//
// go run main.go -cmd=create-application -name="Ionic App"
//
func CreateApplication(db *models.DB, name string) {

	// Generate a random string for the client id.
	clientId, err := helpers.GenerateRandomString(15)

	if err != nil {
		panic(err)
	}

	// Setup the application
	app := models.Application{Name: name, ClientId: clientId, GrantType: "password"}

	// Create new application
	err = db.CreateNewRecord(&app, models.InsertParam{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Application Id: ", app.Id, " ClientId: "+app.ClientId)
}

/* End File */
