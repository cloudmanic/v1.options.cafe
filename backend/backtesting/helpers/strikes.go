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
// GetExpirationDatesFromOptions will return expire dates from an option chain
//
func GetOptionByExpirationDateAndStrike(expire time.Time, strike float64, optionType string, options []types.OptionsChainItem) types.OptionsChainItem {
	rt := types.OptionsChainItem{}

	// Double check TODO(spicer): Return error maybe
	if (optionType != "Put") && (optionType != "Call") {
		return rt
	}

	for _, row := range options {
		if row.OptionType != optionType {
			continue
		}

		if row.ExpirationDate.Format("2006-01-02") != expire.Format("2006-01-02") {
			continue
		}

		if row.Strike != strike {
			continue
		}

		rt = row

		break
	}

	// Return filtered subset
	return rt
}
