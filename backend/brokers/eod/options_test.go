//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-07
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"testing"

	"github.com/nbio/st"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Test - GetGetOptionsChainByExpiration01
//
func TestGetOptionsChainByExpiration01(t *testing.T) {

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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
	st.Expect(t, chain.Puts[50].Strike, 222.00)
	st.Expect(t, chain.Calls[60].Strike, 232.00)
	st.Expect(t, chain.Puts[50].ExpirationDate.Format("2006-01-02"), "2018-10-19")
	st.Expect(t, chain.Calls[60].ExpirationDate.Format("2006-01-02"), "2018-10-19")
}

//
// Test - GetOptionsExpirationsBySymbol01
//
func TestGetOptionsExpirationsBySymbol01(t *testing.T) {

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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
	st.Expect(t, dates[22], "2019-06-21")

}

//
// TestGetOptionsByExpirationType01
//
func TestGetOptionsByExpirationType01(t *testing.T) {

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create broker object
	o := Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// Get options
	options, underlyingLast, err := o.GetOptionsBySymbol("spy")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, underlyingLast, 276.39)
	st.Expect(t, len(options), 6262)

	// Get just PUTs that expire on a certain date.
	list := o.GetOptionsByExpirationType(options[0].ExpirationDate, options[0].OptionType, options)
	st.Expect(t, len(list), 195)
	st.Expect(t, list[0].OptionType, "Call")
	st.Expect(t, list[0].ExpirationDate.Format("2006-01-02"), "2018-10-19")
}

/* End File */
