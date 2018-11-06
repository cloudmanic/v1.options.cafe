//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package paper_trade

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/controllers"
	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/models"
)

//
// Start....
//
func Start() {

	// Lets get started
	services.Critical("Paper Trade Started: " + os.Getenv("APP_ENV"))

	// Start the db connection.
	db, err := models.NewDB()

	if err != nil {
		services.Fatal(err)
	}

	// Startup controller & websockets
	c := &controllers.Controller{DB: db}

	// Start websockets & controllers
	c.StartWebServer()

}

/* End File */
