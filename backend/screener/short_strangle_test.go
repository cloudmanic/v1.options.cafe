//
// Date: 2018-11-08
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-08
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
// Test: RunShortStrangle - 01
//
func TestRunShortStrangle01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "short-strangle",
		Items: []models.ScreenerItem{
			{Key: "put-leg-percent-away", Operator: ">", ValueNumber: 4.5},
			{Key: "call-leg-percent-away", Operator: ">", ValueNumber: 4.5},
			{Key: "open-credit", Operator: ">", ValueNumber: 3.00},
			{Key: "open-credit", Operator: "<", ValueNumber: 3.50},
			{Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// // Setup the broker - Tradier Test
	// broker := tradier.Api{
	//  DB:     nil,
	//  ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	// }

	// Setup the broker - EOD Test
	broker := eod.Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// New screener instance
	s := NewScreen(db, &broker)

	// Run back test
	result, err := s.RunShortStrangle(screen)

	// for _, row := range result {
	// 	fmt.Println(row.Legs[0].OptionStrike, "/", row.Legs[1].OptionStrike, " - ", row.Credit, row.PutPrecentAway, row.CallPrecentAway)
	// }

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, len(result), 17)

	// Result #1
	st.Expect(t, result[0].PutPrecentAway, 5.21)
	st.Expect(t, result[0].CallPrecentAway, 5.29)
	st.Expect(t, result[0].Legs[0].OptionStrike, 262.0)
	st.Expect(t, result[0].Legs[1].OptionStrike, 291.0)
	st.Expect(t, result[0].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[0].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")

	// Result #2
	st.Expect(t, result[10].PutPrecentAway, 5.03)
	st.Expect(t, result[10].CallPrecentAway, 4.92)
	st.Expect(t, result[10].Legs[0].OptionStrike, 262.5)
	st.Expect(t, result[10].Legs[1].OptionStrike, 290.0)
	st.Expect(t, result[10].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[10].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")

	// Result #3
	st.Expect(t, result[13].PutPrecentAway, 5.03)
	st.Expect(t, result[13].CallPrecentAway, 4.74)
	st.Expect(t, result[13].Legs[0].OptionStrike, 262.5)
	st.Expect(t, result[13].Legs[1].OptionStrike, 289.5)
	st.Expect(t, result[13].Legs[0].OptionExpire.Format("2006-01-02"), "2018-11-30")
	st.Expect(t, result[13].Legs[1].OptionExpire.Format("2006-01-02"), "2018-11-30")

}

/* End File */
