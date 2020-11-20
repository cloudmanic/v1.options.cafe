//
// Date: 3/1/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"fmt"
	"reflect"

	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
)

//
// Create models and services in go for Angular
//
// go run main.go -cmd=go-to-angular
//
func GoToAngular() {

	o := models.Settings{}

	t := reflect.TypeOf(o)
	val := reflect.ValueOf(&o).Elem()

	// Build out what we need to map json to Angular vars in the fromJson
	fmt.Println("")
	fmt.Println("########### Json To Angular (used in model::fromJson) ############")
	for i := 0; i < val.NumField(); i++ {
		n := val.Type().Field(i).Name
		f, _ := t.FieldByName(n)

		v, ok := f.Tag.Lookup("json")

		if !ok {
			continue
		}

		fmt.Println("result." + n + ` = json["` + v + `"];`)
	}

	// Build the post back to the server in Angular more or less Angular to json
	fmt.Println("")
	fmt.Println("########### Angular to json (used in service) ############")

	for i := 0; i < val.NumField(); i++ {
		n := val.Type().Field(i).Name
		f, _ := t.FieldByName(n)

		v, ok := f.Tag.Lookup("json")

		if !ok {
			continue
		}

		fmt.Println(v + ": obj." + n + ",")
	}

	fmt.Println("")
}

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
