//
// Date: 2/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nbio/st"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// Test - GetTradeGroups
//
func TestGetTradeGroups01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadTradeGroupTestingData(db.New())

	// Make query
	tgs, _, err := db.GetTradeGroups(QueryParam{PreLoads: []string{"Positions"}})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, tgs[0].Positions[1].Symbol.Id, uint(6))
	st.Expect(t, tgs[0].Positions[1].Symbol.ShortName, "SPY180316P00266000")
	st.Expect(t, tgs[0].Positions[1].Symbol.Name, "SPY Mar 16, 2018 $266.00 Put")
	st.Expect(t, tgs[0].Positions[1].Symbol.Type, "Option")
	st.Expect(t, tgs[0].Positions[1].Symbol.OptionType, "Put")
	st.Expect(t, tgs[0].Positions[1].Symbol.OptionUnderlying, "SPY")
	st.Expect(t, tgs[0].Positions[1].Symbol.OptionStrike, 266.00)
	st.Expect(t, tgs[0].Positions[1].Symbol.OptionExpire.Format("01/02/2006"), "03/16/2018")
}

//
// Test - GetTradeGroupById
//
func TestGetTradeGroupById01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadTradeGroupTestingData(db.New())

	// Make query
	tg, err := db.GetTradeGroupById(1)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, tg.Id, uint(1))
	st.Expect(t, tg.UserId, uint(1))
	st.Expect(t, tg.BrokerAccountId, uint(1))
	st.Expect(t, tg.BrokerAccountRef, "abc123")
	st.Expect(t, tg.Status, "Open")
	st.Expect(t, tg.Type, "Put Credit Spread")
	st.Expect(t, tg.OrderIds, "1")
	st.Expect(t, tg.Risked, 0.00)
	st.Expect(t, tg.PercentGain, 0.00)
	st.Expect(t, tg.Profit, 0.00)
	st.Expect(t, tg.Commission, 23.45)
	st.Expect(t, tg.Note, "Test note #1")

	// Position - 1
	st.Expect(t, len(tg.Positions), 2)
	st.Expect(t, tg.Positions[0].Id, uint(1))
	st.Expect(t, tg.Positions[0].UserId, uint(1))
	st.Expect(t, tg.Positions[0].TradeGroupId, uint(1))
	st.Expect(t, tg.Positions[0].BrokerAccountRef, "123abc")
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
	st.Expect(t, tg.Positions[1].BrokerAccountRef, "123abc")
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

//
// loadTradeGroupTestingData - Load testing data.
//
func loadTradeGroupTestingData(db *gorm.DB) {
	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Users
	db.Create(&User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// Brokers
	db.Create(&Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts})

	// BrokerAccounts
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Add some watchlists
	db.Create(&Watchlist{UserId: 1, Name: "Watchlist #1 - User 1", Symbols: []WatchlistSymbol{{UserId: 1, SymbolId: 1, Order: 0}, {UserId: 1, SymbolId: 2, Order: 1}, {UserId: 1, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 1, Name: "Watchlist #2 - User 1", Symbols: []WatchlistSymbol{{UserId: 1, SymbolId: 1, Order: 0}, {UserId: 1, SymbolId: 2, Order: 1}, {UserId: 1, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 2, Name: "Watchlist #1 - User 2", Symbols: []WatchlistSymbol{{UserId: 2, SymbolId: 1, Order: 0}, {UserId: 2, SymbolId: 2, Order: 1}, {UserId: 2, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 2, Name: "Watchlist #2 - User 2", Symbols: []WatchlistSymbol{{UserId: 2, SymbolId: 1, Order: 0}, {UserId: 2, SymbolId: 2, Order: 1}, {UserId: 2, SymbolId: 3, Order: 2}}})

	// TradeGroups TODO: make this more complete
	db.Create(&TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Put Credit Spread", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})

	// Positions TODO: Put better values in here.
	db.Create(&Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	db.Create(&Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})
}

/* End File */
