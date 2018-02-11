//
// Date: 2/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"

	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test - GetTradeGroupById
//
func TestGetTradeGroupById01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Make query
	tg, err := db.GetTradeGroupById(1)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, tg.Id, uint(1))
	st.Expect(t, tg.UserId, uint(1))
	st.Expect(t, tg.BrokerAccountId, uint(1))
	st.Expect(t, tg.AccountId, "abc123")
	st.Expect(t, tg.Status, "Open")
	st.Expect(t, tg.Type, "Put Credit Spread")
	st.Expect(t, tg.OrderIds, "1")
	st.Expect(t, tg.Risked, 0.00)
	st.Expect(t, tg.Gain, 0.00)
	st.Expect(t, tg.Profit, 0.00)
	st.Expect(t, tg.Commission, 23.45)
	st.Expect(t, tg.Note, "Test note #1")

	// Position - 1
	st.Expect(t, len(tg.Positions), 2)
	st.Expect(t, tg.Positions[0].Id, uint(1))
	st.Expect(t, tg.Positions[0].UserId, uint(1))
	st.Expect(t, tg.Positions[0].TradeGroupId, uint(1))
	st.Expect(t, tg.Positions[0].AccountId, "123abc")
	st.Expect(t, tg.Positions[0].Status, "Open")
	st.Expect(t, tg.Positions[0].SymbolId, uint(4))
	st.Expect(t, tg.Positions[0].Qty, 10)
	st.Expect(t, tg.Positions[0].OrgQty, 10)
	st.Expect(t, tg.Positions[0].CostBasis, 1000.00)
	st.Expect(t, tg.Positions[0].AvgOpenPrice, 1.00)
	st.Expect(t, tg.Positions[0].AvgClosePrice, 0.00)
	st.Expect(t, tg.Positions[0].OrderIds, "1")
	st.Expect(t, tg.Positions[0].Note, "Test note #1")

	// Symnbol - 1
	st.Expect(t, tg.Positions[0].Symbol.Id, uint(4))
	st.Expect(t, tg.Positions[0].Symbol.ShortName, "SPY180316P00253000")
	st.Expect(t, tg.Positions[0].Symbol.Name, "SPY Mar 16, 2018 $253.00 Put")
	st.Expect(t, tg.Positions[0].Symbol.Type, "Option")

	// Position - 2
	st.Expect(t, tg.Positions[1].Id, uint(2))
	st.Expect(t, tg.Positions[1].UserId, uint(1))
	st.Expect(t, tg.Positions[1].TradeGroupId, uint(1))
	st.Expect(t, tg.Positions[1].AccountId, "123abc")
	st.Expect(t, tg.Positions[1].Status, "Open")
	st.Expect(t, tg.Positions[1].SymbolId, uint(6))
	st.Expect(t, tg.Positions[1].Qty, 10)
	st.Expect(t, tg.Positions[1].OrgQty, 10)
	st.Expect(t, tg.Positions[1].CostBasis, 1000.00)
	st.Expect(t, tg.Positions[1].AvgOpenPrice, 1.00)
	st.Expect(t, tg.Positions[1].AvgClosePrice, 0.00)
	st.Expect(t, tg.Positions[1].OrderIds, "1")
	st.Expect(t, tg.Positions[1].Note, "Test note #2")

	// Symnbol - 2
	st.Expect(t, tg.Positions[1].Symbol.Id, uint(6))
	st.Expect(t, tg.Positions[1].Symbol.ShortName, "SPY180316P00266000")
	st.Expect(t, tg.Positions[1].Symbol.Name, "SPY Mar 16, 2018 $266.00 Put")
	st.Expect(t, tg.Positions[1].Symbol.Type, "Option")
}

/* End File */
