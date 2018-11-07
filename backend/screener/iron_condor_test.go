//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-07
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// RunIronCondor - 01
//
func TestRunIronCondor01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "iron-condor",
		Items: []models.ScreenerItem{
			{Key: "put-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "call-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "put-leg-percent-away", Operator: ">", ValueNumber: 4.5},
			{Key: "call-leg-percent-away", Operator: ">", ValueNumber: 4.5},
			{Key: "open-credit", Operator: ">", ValueNumber: 0.50},
			{Key: "open-credit", Operator: "<", ValueNumber: 3.00},
			{Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// // Setup the broker - Tradier Test
	// broker := tradier.Api{
	// 	DB:     nil,
	// 	ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	// }

	// Setup the broker - EOD Test
	broker := eod.Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// New screener instance
	s := NewScreen(db, &broker)

	// Run back test
	result, err := s.RunIronCondor(screen)

	// Test result - ORDER seems to be different each time
	st.Expect(t, err, nil)
	st.Expect(t, len(result), 3)

	// Result #1
	st.Expect(t, result[0].PutPrecentAway, 4.84)
	st.Expect(t, result[0].CallPrecentAway, 4.69)
	st.Expect(t, result[0].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[0].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[0].Legs[2].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[0].Legs[3].OptionExpire.Format("2006-01-02"), "2018-11-30")

	// Result #2
	st.Expect(t, result[1].PutPrecentAway, 5.21)
	st.Expect(t, result[1].CallPrecentAway, 4.36)
	st.Expect(t, result[1].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[1].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[1].Legs[2].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[1].Legs[3].OptionExpire.Format("2006-01-02"), "2018-11-30")

	// Result #3
	st.Expect(t, result[2].PutPrecentAway, 4.84)
	st.Expect(t, result[2].CallPrecentAway, 4.36)
	st.Expect(t, result[2].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[2].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[2].Legs[2].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[2].Legs[3].OptionExpire.Format("2006-01-02"), "2018-11-30")

}

/* End File */
