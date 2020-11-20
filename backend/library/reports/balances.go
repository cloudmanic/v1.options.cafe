//
// Date: 3/78/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"time"

	"app.options.cafe/models"
	"github.com/jinzhu/now"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

// BalancesParams struct
type BalancesParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GroupBy   string    `json:group_by`
	Sort      string    `json:sort`
}

// Balance struct
type Balances struct {
	Date         Date    `json:"date"`
	DateStr      string  `json:"-"`
	TotalCash    float64 `json:"total_cash"`
	AccountValue float64 `json:"account_value"`
}

// Returns struct
type Returns struct {
	Date         Date    `json:"date"`
	Percent      float64 `json:"percent"`
	AccountValue float64 `json:"account_value"`
	TotalCash    float64 `json:"total_cash"`
	PricePer     float64 `json:"price_per"`
	Units        float64 `json:"units"`
}

//
// GetBalances based on parms we pass in.
//
func GetBalances(db models.Datastore, brokerAccount models.BrokerAccount, parms BalancesParams) []Balances {
	selectStr := ""
	queryStr := ""
	balances := []Balances{}

	// Build query
	switch parms.GroupBy {

	case "day":
		selectStr = `DATE_FORMAT(date,'%Y-%m-%d') AS date_str, total_cash, account_value`
		queryStr = "SELECT " + selectStr + " FROM balance_histories WHERE broker_account_id = ? AND date >= ? AND date <= ? GROUP BY date " + parms.Sort

		// Query
		db.New().Raw(queryStr, brokerAccount.Id, parms.StartDate, parms.EndDate).Scan(&balances)

	case "month":
		wIn := `id IN ( SELECT MAX(id) FROM balance_histories WHERE (broker_account_id = ? AND date >= ? AND date <= ?) GROUP BY DATE_FORMAT(date,'%Y-%m') )`
		selectStr = `DATE_FORMAT(date,'%Y-%m') AS date_str, total_cash, account_value`
		queryStr = "SELECT " + selectStr + " FROM balance_histories WHERE broker_account_id = ? AND date >= ? AND date <= ? AND " + wIn + " GROUP BY DATE_FORMAT(date,'%Y-%m') " + parms.Sort

		// Query
		db.New().Raw(queryStr, brokerAccount.Id, parms.StartDate, parms.EndDate, brokerAccount.Id, parms.StartDate, parms.EndDate).Scan(&balances)

	case "year":
		wIn := `id IN ( SELECT MAX(id) FROM balance_histories WHERE (broker_account_id = ? AND date >= ? AND date <= ?) GROUP BY YEAR(date) )`
		selectStr = `YEAR(date) AS date_str, total_cash, account_value`
		queryStr = "SELECT " + selectStr + " FROM balance_histories WHERE broker_account_id = ? AND date >= ? AND date <= ? AND " + wIn + " GROUP BY YEAR(date) " + parms.Sort

		// Query
		db.New().Raw(queryStr, brokerAccount.Id, parms.StartDate, parms.EndDate, brokerAccount.Id, parms.StartDate, parms.EndDate).Scan(&balances)
	}

	// Post query processing
	for key, row := range balances {
		balances[key].Date = Date{now.New(helpers.ParseDateNoError(row.DateStr)).EndOfMonth()}

		// Special case for year grouping
		if parms.GroupBy == "year" {
			balances[key].Date = Date{now.New(helpers.ParseDateNoError(row.DateStr + "-12")).EndOfMonth()}
		}

		// Special case for day grouping
		if parms.GroupBy == "day" {
			balances[key].Date = Date{helpers.ParseDateNoError(row.DateStr)}
		}
	}

	// Return Happy
	return balances
}

//
// GetAccountReturns - return returns based on per shares
//
func GetAccountReturns(db models.Datastore, brokerAccount models.BrokerAccount, parms BalancesParams) []Returns {
	returns := []Returns{}
	achMap := make(map[string]float64)

	// We want to build a profit of money in and money out of the system
	be := []models.BrokerEvent{}
	db.New().Where("description = ? OR description = ? OR description = ?", "ACH DEPOSIT", "ACH DIRECT WITHDRAWAL", "ACH DISBURSEMENT").Where("type = ? OR type = ?", "Journal", "Ach").Where("broker_account_id = ? AND date >= ? AND date <= ?", brokerAccount.Id, parms.StartDate, parms.EndDate).Order("date ASC").Find(&be)

	for _, row := range be {
		achMap[row.Date.Format("2006-01-02")] = row.Amount
	}

	// Get balanaces
	parms.GroupBy = "day"
	parms.Sort = "asc"
	balances := GetBalances(db, brokerAccount, parms)

	// Make sure we have at least one balance
	if len(balances) <= 0 {
		return returns
	}

	// Set starting units
	units := balances[0].AccountValue

	// Loop through the balances and set returns
	for key, row := range balances {
		// Set price per
		pricePer := row.AccountValue / units

		// See if we had any ACHs today. Key > 0 solves for case when start day we also added money.
		if key > 0 {
			if val, ok := achMap[row.Date.Format("2006-01-02")]; ok {

				// Update the price per.
				pricePer = (row.AccountValue - val) / units

				// Adding or removing money to account
				if val > 0 {
					units = units + (val / pricePer)
				} else if val < 0 {
					units = units - (val / pricePer)
				}
			}
		}

		// Add onto returns array
		returns = append(returns, Returns{
			Date:         row.Date,
			Percent:      (pricePer - 1.00),
			PricePer:     pricePer,
			Units:        units,
			AccountValue: row.AccountValue,
			TotalCash:    row.TotalCash,
		})
	}

	// Return happy
	return returns
}

/* End File */
