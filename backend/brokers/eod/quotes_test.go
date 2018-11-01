//
// Date: 2018-11-01
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-01
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package eod

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test - GetQuotes - 01
//
func TestGetQuotes01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create broker object
	broker := Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// Get quote from S3 store
	quotes, err := broker.GetQuotes([]string{"spy", "iwm"})

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, len(quotes), 2)
	st.Expect(t, quotes[0].Last, 276.39)
	st.Expect(t, quotes[1].Last, 155.01)
	st.Expect(t, quotes[0].Symbol, "SPY")
	st.Expect(t, quotes[1].Symbol, "IWM")
}

/* End File */
