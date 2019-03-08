//
// Date: 3/78/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
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

/* End File */
