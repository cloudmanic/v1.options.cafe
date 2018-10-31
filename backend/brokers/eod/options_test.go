//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/nbio/st"
)

//
// Test - GetOptionsExpirationsBySymbol01
//
func TestGetOptionsExpirationsBySymbol01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create broker object
	o := Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// Get dates from S3 store
	dates, err := o.GetOptionsExpirationsBySymbol("spy")

	spew.Dump(dates)

	// Test result
	st.Expect(t, err, nil)
	//st.Expect(t, (len(dates) >= 3477), true)

}

/* End File */
