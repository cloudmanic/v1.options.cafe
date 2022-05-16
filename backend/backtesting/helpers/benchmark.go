//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"time"

	"app.options.cafe/brokers/types"
)

//
// GetBenchmarkByDate will return the close price of the benchmark
//
func GetBenchmarkByDate(date time.Time, benchmarkQuotes []types.HistoryQuote) float64 {
	for _, row := range benchmarkQuotes {
		if row.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			return row.Close
		}
	}
	return 0.00
}
