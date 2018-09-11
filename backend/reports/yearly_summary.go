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
	Year           int     `json:"year"`
	TotalTrades    int     `json:"total_trades"`
	LossCount      int     `json:"loss_count"`
	WinCount       int     `json:"win_count"`
	Profit         float64 `json:"profit"`
	Commission     float64 `json:"commission"`
	WinPercent     float64 `json:"win_percent"`
	LossPercent    float64 `json:"loss_percent"`
	ProfitStd      float64 `json:"profit_std"`
	PercentGainStd float64 `json:"precent_gain_std"`
	SharpeRatio    float64 `json:"sharpe_ratio"`
	AvgRisked      float64 `json:"avg_risked"`
	AvgPercentGain float64 `json:"avg_percent_gain"`
}

//
// Return a list of years we have trade groups for.
//
func GetYearsWithTradeGroups(db models.Datastore, brokerAccount models.BrokerAccount) []int {

	years := []int{}

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

	if summary.TotalTrades > 0 {
		summary.WinPercent = helpers.Round((float64(summary.WinCount)/float64(summary.TotalTrades))*100, 2)
		summary.LossPercent = helpers.Round((float64(summary.LossCount)/float64(summary.TotalTrades))*100, 2)
		summary.SharpeRatio = helpers.Round((summary.AvgPercentGain-riskFreeRate)/summary.PercentGainStd, 2)
		summary.ProfitStd = helpers.Round(summary.ProfitStd, 2)
		summary.AvgRisked = helpers.Round(summary.AvgRisked, 2)
		summary.PercentGainStd = helpers.Round(summary.PercentGainStd, 2)
	}

	return summary
}

/* End File */
