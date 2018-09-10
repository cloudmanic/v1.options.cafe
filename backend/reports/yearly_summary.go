//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type YearlySummary struct {
	Year           int
	TotalTrades    int
	LossCount      int
	WinCount       int
	Profit         float64
	Commission     float64
	WinPercent     float64
	LossPercent    float64
	ProfitStd      float64
	PercentGainStd float64
	SharpeRatio    float64
	AvgRisked      float64
	AvgPercentGain float64
}

//
// Return a list of years we have trade groups for.
//
func GetYearsWithTradeGroups(db models.Datastore, brokerAccount models.BrokerAccount) []int {

	var years []int

	type Result struct {
		Year int
	}

	var results []Result

	// Set select string
	queryStr := "SELECT YEAR(closed_date) AS year FROM trade_groups WHERE broker_account_id = ? GROUP BY year ORDER BY year desc"

	// Run query.
	db.New().Raw(queryStr, brokerAccount.Id).Scan(&results)

	for _, row := range results {
		years = append(years, row.Year)
	}

	return years
}

//
// Get a yearly summary based on account, year
//
func GetYearlySummaryByAccountYear(db models.Datastore, brokerAccount models.BrokerAccount, year int) YearlySummary {

	riskFreeRate := 2.03 // TODO: automate setting this
	summary := YearlySummary{}

	// Build query
	selectStr := "COUNT(profit) AS total_trades, SUM(profit) AS profit, SUM(commission) AS commission, " +
		"SUM(CASE WHEN profit < 0 THEN 1 ELSE 0 END) AS loss_count, SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) AS win_count, " +
		"STDDEV(profit) AS profit_std, AVG(percent_gain) AS avg_percent_gain, AVG(risked) AS avg_risked, STDDEV(percent_gain) AS percent_gain_std"

	queryStr := "SELECT " + selectStr + " FROM trade_groups WHERE broker_account_id = ? AND YEAR(closed_date) = ?"

	// Run query.
	db.New().Raw(queryStr, brokerAccount.Id, year).Scan(&summary)

	// Post query processing
	summary.Year = year
	summary.WinPercent = helpers.Round((float64(summary.WinCount)/float64(summary.TotalTrades))*100, 2)
	summary.LossPercent = helpers.Round((float64(summary.LossCount)/float64(summary.TotalTrades))*100, 2)
	summary.SharpeRatio = helpers.Round((summary.AvgPercentGain-riskFreeRate)/summary.PercentGainStd, 2)
	summary.ProfitStd = helpers.Round(summary.ProfitStd, 2)
	summary.AvgRisked = helpers.Round(summary.AvgRisked, 2)
	summary.PercentGainStd = helpers.Round(summary.PercentGainStd, 2)

	return summary
}

/* End File */
