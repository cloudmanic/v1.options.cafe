//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/jinzhu/now"
)

type ProfitLossParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GroupBy   string    `json:group_by`
	Sort      string    `json:sort`
}

type ProfitLoss struct {
	Date           Date    `json:"date"`
	DateStr        string  `json:"-"`
	Profit         float64 `json:"profit"`
	TradeCount     int     `json:"trade_count"`
	Commissions    float64 `json:"commissions"`
	ProfitPerTrade float64 `json:"profit_per_trade"`
	WinRatio       float64 `json:"win_ratio"`
	LossCount      int     `json:"loss_count"`
	WinCount       int     `json:"win_count"`
}

//
// Get profit and loss based on parms we pass in.
//
func GetProfitLoss(db models.Datastore, brokerAccount models.BrokerAccount, parms ProfitLossParams) []ProfitLoss {

	selectStr := ""
	queryStr := ""
	profits := []ProfitLoss{}

	// Build query
	switch parms.GroupBy {

	case "day":
		selectStr = `DATE_FORMAT(closed_date,'%Y-%m-%d') AS date_str, SUM(profit) AS profit, SUM(commission) AS commissions, COUNT(closed_date) AS trade_count, SUM(CASE WHEN profit < 0 THEN 1 ELSE 0 END) AS loss_count, SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) AS win_count`
		queryStr = "SELECT " + selectStr + " FROM trade_groups WHERE status = 'Closed' AND broker_account_id = ? AND open_date >= ? AND closed_date <= ? GROUP BY DATE_FORMAT(closed_date,'%Y-%m-%d') " + parms.Sort

	case "month":
		selectStr = `DATE_FORMAT(closed_date,'%Y-%m') AS date_str, SUM(profit) AS profit, SUM(commission) AS commissions, COUNT(closed_date) AS trade_count, SUM(CASE WHEN profit < 0 THEN 1 ELSE 0 END) AS loss_count, SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) AS win_count`
		queryStr = "SELECT " + selectStr + " FROM trade_groups WHERE status = 'Closed' AND broker_account_id = ? AND open_date >= ? AND closed_date <= ? GROUP BY DATE_FORMAT(closed_date,'%Y-%m') " + parms.Sort

	case "year":
		selectStr = `YEAR(closed_date) AS date_str, SUM(profit) AS profit, SUM(commission) AS commissions, COUNT(closed_date) AS trade_count, SUM(CASE WHEN profit < 0 THEN 1 ELSE 0 END) AS loss_count, SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) AS win_count`
		queryStr = "SELECT " + selectStr + " FROM trade_groups WHERE status = 'Closed' AND broker_account_id = ? AND open_date >= ? AND closed_date <= ? GROUP BY YEAR(closed_date) " + parms.Sort

	}

	// Run query.
	db.New().Raw(queryStr, brokerAccount.Id, parms.StartDate, parms.EndDate).Scan(&profits)

	// Post query processing
	for key, row := range profits {
		profits[key].WinRatio = (float64(row.WinCount) / float64(row.TradeCount)) * 100
		profits[key].ProfitPerTrade = row.Profit / float64(row.TradeCount)
		profits[key].Date = Date{now.New(helpers.ParseDateNoError(row.DateStr)).EndOfMonth()}

		// Special case for year grouping
		if parms.GroupBy == "year" {
			profits[key].Date = Date{now.New(helpers.ParseDateNoError(row.DateStr + "-12")).EndOfMonth()}
		}

	}

	// Return Happy
	return profits
}

/* End File */
