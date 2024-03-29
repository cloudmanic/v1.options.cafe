//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package market

import (
	"testing"
)

//
// Test - GetMarketStatus
//
func TestGetMarketStatus01(t *testing.T) {

	// // Flush pending mocks after test execution
	// defer gock.Off()

	// // Setup mock request.
	// gock.New("https://api.tradier.com/v1").
	// 	Get("/markets/clock").
	// 	Reply(200).
	// 	BodyString(`{"clock":{"state":"closed","date":"2017-09-04","timestamp":1504505290,"next_state":"premarket","next_change":"08:00","description":"Market is closed."}}`)

	// // Make API call
	// marketStatus, err := GetMarketStatus()

	// if err != nil {
	// 	panic(err)
	// }

	// // Verify the data was return as expected
	// st.Expect(t, marketStatus.Date, "2017-09-04")
	// st.Expect(t, marketStatus.State, "closed")
	// st.Expect(t, marketStatus.Description, "Market is closed.")

	// // Verify that we don't have pending mocks
	// st.Expect(t, gock.IsDone(), true)
}

/* End File */
