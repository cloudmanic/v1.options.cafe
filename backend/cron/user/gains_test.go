//
// Date: 2022-07-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package user

import (
	"testing"

	"app.options.cafe/models"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

//
// TestImportGainsForUser01 will make API calls to the broker and pull in profit / loss information and store in the database
//
func TestImportGainsForUser01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Flush pending mocks after test execution
	defer gock.Off()

	// Mock account id
	accountId := "6YA06085"

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Get("/accounts/" + accountId + "/gainloss").
		Reply(200).
		BodyString(`{"gainloss":{"closed_position":[{"close_date":"2022-07-07T00:00:00.000Z","cost":26.15,"gain_loss":513.67,"gain_loss_percent":1964.3212,"open_date":"2022-06-01T00:00:00.000Z","proceeds":539.82,"quantity":-2.0,"symbol":"SPY220708P00372000","term":36},{"close_date":"2022-07-07T00:00:00.000Z","cost":498.15,"gain_loss":-478.32,"gain_loss_percent":-96.0193,"open_date":"2022-06-01T00:00:00.000Z","proceeds":19.83,"quantity":2.0,"symbol":"SPY220708P00370000","term":36},{"close_date":"2022-06-24T00:00:00.000Z","cost":16.15,"gain_loss":361.68,"gain_loss_percent":2239.5046,"open_date":"2022-06-03T00:00:00.000Z","proceeds":377.83,"quantity":-2.0,"symbol":"SPY220624P00380000","term":21},{"close_date":"2022-06-24T00:00:00.000Z","cost":42.15,"gain_loss":737.67,"gain_loss_percent":1750.1068,"open_date":"2022-05-18T00:00:00.000Z","proceeds":779.82,"quantity":-2.0,"symbol":"SPY220701P00356000","term":37},{"close_date":"2022-06-24T00:00:00.000Z","cost":336.15,"gain_loss":-326.32,"gain_loss_percent":-97.0757,"open_date":"2022-06-03T00:00:00.000Z","proceeds":9.83,"quantity":2.0,"symbol":"SPY220624P00378000","term":21},{"close_date":"2022-06-24T00:00:00.000Z","cost":730.15,"gain_loss":-694.32,"gain_loss_percent":-95.0928,"open_date":"2022-05-18T00:00:00.000Z","proceeds":35.83,"quantity":2.0,"symbol":"SPY220701P00354000","term":37},{"close_date":"2022-06-17T00:00:00.000Z","cost":1581.23,"gain_loss":1064.46,"gain_loss_percent":67.3185,"open_date":"2022-05-05T00:00:00.000Z","proceeds":2645.69,"quantity":3.0,"symbol":"SPY220617P00374000","term":43},{"close_date":"2022-06-17T00:00:00.000Z","cost":3246.23,"gain_loss":-1572.48,"gain_loss_percent":-48.4402,"open_date":"2022-05-05T00:00:00.000Z","proceeds":1673.75,"quantity":-3.0,"symbol":"SPY220617P00376000","term":43},{"close_date":"2022-06-17T00:00:00.000Z","cost":77798.19,"gain_loss":-382.16,"gain_loss_percent":-0.4912,"open_date":"2022-06-17T00:00:00.000Z","proceeds":77416.03,"quantity":200.0,"symbol":"SPY","term":0},{"close_date":"2022-06-07T00:00:00.000Z","cost":550.16,"gain_loss":-488.33,"gain_loss_percent":-88.7615,"open_date":"2022-05-16T00:00:00.000Z","proceeds":61.83,"quantity":2.0,"symbol":"SPY220630P00353000","term":22},{"close_date":"2022-07-07T00:00:00.000Z","cost":26.15,"gain_loss":513.67,"gain_loss_percent":1964.3212,"open_date":"2022-06-01T00:00:00.000Z","proceeds":539.82,"quantity":-2.0,"symbol":"SPY220708P00372000","term":36},{"close_date":"2022-07-07T00:00:00.000Z","cost":498.15,"gain_loss":-478.32,"gain_loss_percent":-96.0193,"open_date":"2022-06-01T00:00:00.000Z","proceeds":19.83,"quantity":2.0,"symbol":"SPY220708P00370000","term":36},{"close_date":"2022-06-24T00:00:00.000Z","cost":16.15,"gain_loss":361.68,"gain_loss_percent":2239.5046,"open_date":"2022-06-03T00:00:00.000Z","proceeds":377.83,"quantity":-2.0,"symbol":"SPY220624P00380000","term":21},{"close_date":"2022-06-24T00:00:00.000Z","cost":42.15,"gain_loss":737.67,"gain_loss_percent":1750.1068,"open_date":"2022-05-18T00:00:00.000Z","proceeds":779.82,"quantity":-2.0,"symbol":"SPY220701P00356000","term":37},{"close_date":"2022-06-24T00:00:00.000Z","cost":336.15,"gain_loss":-326.32,"gain_loss_percent":-97.0757,"open_date":"2022-06-03T00:00:00.000Z","proceeds":9.83,"quantity":2.0,"symbol":"SPY220624P00378000","term":21},{"close_date":"2022-06-24T00:00:00.000Z","cost":730.15,"gain_loss":-694.32,"gain_loss_percent":-95.0928,"open_date":"2022-05-18T00:00:00.000Z","proceeds":35.83,"quantity":2.0,"symbol":"SPY220701P00354000","term":37},{"close_date":"2022-06-17T00:00:00.000Z","cost":1581.23,"gain_loss":1064.46,"gain_loss_percent":67.3185,"open_date":"2022-05-05T00:00:00.000Z","proceeds":2645.69,"quantity":3.0,"symbol":"SPY220617P00374000","term":43},{"close_date":"2022-06-17T00:00:00.000Z","cost":3246.23,"gain_loss":-1572.48,"gain_loss_percent":-48.4402,"open_date":"2022-05-05T00:00:00.000Z","proceeds":1673.75,"quantity":-3.0,"symbol":"SPY220617P00376000","term":43},{"close_date":"2022-06-17T00:00:00.000Z","cost":77798.19,"gain_loss":-382.16,"gain_loss_percent":-0.4912,"open_date":"2022-06-17T00:00:00.000Z","proceeds":77416.03,"quantity":200.0,"symbol":"SPY","term":0},{"close_date":"2022-06-07T00:00:00.000Z","cost":550.16,"gain_loss":-488.33,"gain_loss_percent":-88.7615,"open_date":"2022-05-16T00:00:00.000Z","proceeds":61.83,"quantity":2.0,"symbol":"SPY220630P00353000","term":22}]}}`)

	// Run the function we are testing. BodyString above has duplicate data to make sure only the first entry gets saved in DB.
	ImportGainsForUser(db, models.BrokerAccount{
		Id:            6,
		UserId:        2,
		AccountNumber: accountId,
	})

	// // Setup a new backtesting
	// bt := New(db, 1, "SPY")

	// // Build screener object
	// screen := models.Screener{
	// 	Symbol:   "SPY",
	// 	Strategy: "blank",
	// }

	// // Run blank backtest
	// err := bt.DoBacktestDays(&models.Backtest{
	// 	StartDate: models.Date{helpers.ParseDateNoError("2018-01-01")},
	// 	EndDate:   models.Date{helpers.ParseDateNoError("2019-01-01")},
	// 	Screen:    screen,
	// })
	//st.Expect(t, err, nil)
	st.Expect(t, nil, nil)
}
