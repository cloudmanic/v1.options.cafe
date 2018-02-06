//
// Date: 2/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"testing"

	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - GetPositions - More than one account.
//
func TestGetPositions01(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/user/positions").
		Reply(200).
		BodyString(`{"accounts":{"account":[{"account_number":"6YA06984","positions":"null"},{"account_number":"6YA05782","positions":{"position":[{"cost_basis":519.5,"date_acquired":"2018-01-09T11:15:51.981Z","id":266353,"quantity":10.0,"symbol":"CONE"},{"cost_basis":387.2796,"date_acquired":"2018-01-09T11:15:51.999Z","id":266354,"quantity":4.0,"symbol":"DIS"},{"cost_basis":395.7996,"date_acquired":"2018-01-09T11:15:52.016Z","id":266355,"quantity":4.0,"symbol":"NFLX"},{"cost_basis":169.598,"date_acquired":"2018-01-09T11:15:52.034Z","id":266356,"quantity":20.0,"symbol":"USO"},{"cost_basis":484.8396,"date_acquired":"2018-01-09T11:15:52.052Z","id":266357,"quantity":4.0,"symbol":"VTI"}]}},{"account_number":"6YA06085","positions":{"position":[{"cost_basis":1404.0,"date_acquired":"2018-02-05T16:17:18.308Z","id":286040,"quantity":9.0,"symbol":"SPY180309P00262000"},{"cost_basis":-1629.0,"date_acquired":"2018-02-05T16:17:18.308Z","id":286039,"quantity":-9.0,"symbol":"SPY180309P00264000"},{"cost_basis":2303.01,"date_acquired":"2018-02-05T20:14:14.167Z","id":286459,"quantity":9.0,"symbol":"SPY180316P00250000"},{"cost_basis":2097.0,"date_acquired":"2018-02-05T19:26:06.234Z","id":286357,"quantity":9.0,"symbol":"SPY180316P00251000"},{"cost_basis":-2618.01,"date_acquired":"2018-02-05T20:14:14.167Z","id":286458,"quantity":-9.0,"symbol":"SPY180316P00252000"},{"cost_basis":-2295.0,"date_acquired":"2018-02-05T19:26:06.234Z","id":286356,"quantity":-9.0,"symbol":"SPY180316P00253000"},{"cost_basis":1332.0,"date_acquired":"2018-02-02T16:19:15.693Z","id":284951,"quantity":9.0,"symbol":"SPY180316P00264000"},{"cost_basis":-1539.0,"date_acquired":"2018-02-02T16:19:15.693Z","id":284952,"quantity":-9.0,"symbol":"SPY180316P00266000"},{"cost_basis":-1270.0,"date_acquired":"2018-02-06T17:39:07.118Z","id":286943,"quantity":-2.0,"symbol":"VXX180223C00050000"},{"cost_basis":982.0,"date_acquired":"2018-02-06T17:39:07.118Z","id":286944,"quantity":2.0,"symbol":"VXX180223C00055000"},{"cost_basis":1290.0,"date_acquired":"2018-02-06T17:36:58.677Z","id":286939,"quantity":2.0,"symbol":"VXX180302P00046000"}]}}]}}`)

	// Create new tradier instance
	tradier := &Api{}

	// Make API call
	positions, err := tradier.GetPositions()

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, positions[0].Id, 266353)
	st.Expect(t, positions[0].AccountId, "6YA05782")
	st.Expect(t, positions[0].Symbol, "CONE")
	st.Expect(t, positions[0].DateAcquired, "2018-01-09T11:15:51.981Z")
	st.Expect(t, positions[0].CostBasis, 519.5)
	st.Expect(t, positions[0].Quantity, float64(10))

	st.Expect(t, positions[1].Id, 266354)
	st.Expect(t, positions[1].AccountId, "6YA05782")
	st.Expect(t, positions[1].Symbol, "DIS")
	st.Expect(t, positions[1].DateAcquired, "2018-01-09T11:15:51.999Z")
	st.Expect(t, positions[1].CostBasis, 387.2796)
	st.Expect(t, positions[1].Quantity, float64(4))

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

//
// Test - GetPositions - Just one account
//
func TestGetPositions02(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/user/positions").
		Reply(200).
		BodyString(`{"accounts":{"account":{"account_number":"6YA05782","positions":{"position":[{"cost_basis":519.5,"date_acquired":"2018-01-09T11:15:51.981Z","id":266353,"quantity":10,"symbol":"CONE"},{"cost_basis":387.2796,"date_acquired":"2018-01-09T11:15:51.999Z","id":266354,"quantity":4,"symbol":"DIS"},{"cost_basis":395.7996,"date_acquired":"2018-01-09T11:15:52.016Z","id":266355,"quantity":4,"symbol":"NFLX"},{"cost_basis":169.598,"date_acquired":"2018-01-09T11:15:52.034Z","id":266356,"quantity":20,"symbol":"USO"},{"cost_basis":484.8396,"date_acquired":"2018-01-09T11:15:52.052Z","id":266357,"quantity":4,"symbol":"VTI"}]}}}}`)
	// Create new tradier instance
	tradier := &Api{}

	// Make API call
	positions, err := tradier.GetPositions()

	if err != nil {
		panic(err)
	}

	// // Verify the data was return as expected
	st.Expect(t, positions[0].Id, 266353)
	st.Expect(t, positions[0].AccountId, "6YA05782")
	st.Expect(t, positions[0].Symbol, "CONE")
	st.Expect(t, positions[0].DateAcquired, "2018-01-09T11:15:51.981Z")
	st.Expect(t, positions[0].CostBasis, 519.5)
	st.Expect(t, positions[0].Quantity, float64(10))

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

/* End File */
