//
// Date: 2018-12-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package market

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - Return a underlying stock quote based on date in time.
//
func TestGetUnderlayingQuoteByDate01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Clear tables - Default test values throws API Token to Tradier off.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Get("/markets/history").
		Reply(200).
		BodyString(`{"history":{"day":{"date":"2018-12-21","open":246.74,"high":249.71,"low":239.98,"close":240.7,"volume":255320360}}}`)

	// Make API call
	quote, err := GetUnderlayingQuoteByDate(db, uint(1), "spy", helpers.ParseDateNoError("12/21/18"))

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, quote, 240.70)

}

/* End File */
