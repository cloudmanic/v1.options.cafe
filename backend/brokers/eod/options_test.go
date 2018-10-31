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

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test - GetGetOptionsChainByExpiration01
//
func TestGetOptionsChainByExpiration01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create broker object
	o := Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// Get chain from S3 store
	chain, err := o.GetOptionsChainByExpiration("spy", "2018-10-19")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, len(chain.Calls), 195)
	st.Expect(t, len(chain.Puts), 195)
	st.Expect(t, chain.Puts[50].ExpirationDate, types.Date{helpers.ParseDateNoError("2018-10-19").UTC()})
	st.Expect(t, chain.Calls[60].ExpirationDate, types.Date{helpers.ParseDateNoError("2018-10-19").UTC()})
}

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

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, len(dates), 33)

}

/* End File */
