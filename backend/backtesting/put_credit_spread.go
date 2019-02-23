//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
)

//
// DoPutCreditSpread - Run a put credit spread backtest.
//
func (t *Base) DoPutCreditSpread(today time.Time, action Action, chains map[time.Time]types.OptionsChain) error {
	// jsonRt, _ := json.Marshal(options)
	//
	// fmt.Println(len(string(jsonRt)))
	// os.Exit(1)

	// (types.OptionsChainItem) {
	//  Underlying: (string) (len=3) "SPY",
	//  Symbol: (string) "",
	//  OptionType: (string) (len=4) "Call",
	//  Description: (string) (len=29) "SPY Oct 19, 2018 $100.00 Call",
	//  Strike: (float64) 100,
	//  ExpirationDate: (types.Date) 2018-10-19 07:00:00 +0000 UTC,
	//  Last: (float64) 177.57,
	//  Change: (float64) 0,
	//  ChangePercentage: (float64) 0,
	//  Volume: (int) 0,
	//  AverageVolume: (int) 0,
	//  LastVolume: (int) 0,
	//  Open: (float64) 0,
	//  High: (float64) 0,
	//  Low: (float64) 0,
	//  Close: (float64) 0,
	//  Bid: (float64) 176.61,
	//  BidSize: (int) 0,
	//  Ask: (float64) 177.16,
	//  AskSize: (int) 0,
	//  OpenInterest: (int) 775,
	//  ImpliedVol: (float64) 177.16,
	//  Delta: (float64) 0.9871,
	//  Gamma: (float64) 0.0233,
	//  Theta: (float64) -238.8362,
	//  Vega: (float64) 0.4798
	// }

	//fmt.Println(underlyingLast, " : ", row.Format("2006-01-02"), " : ", len(options))

	return nil
}

/* End File */
