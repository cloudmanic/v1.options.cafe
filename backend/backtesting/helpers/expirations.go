//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import "app.options.cafe/brokers/types"

//
// GetExpirationDatesFromOptions will return expire dates from an option chain
//
func GetExpirationDatesFromOptions(options []types.OptionsChainItem) []types.Date {
	seen := map[types.Date]bool{}
	dates := []types.Date{}

	// Loop through and get expire dates
	for _, row := range options {
		if _, ok := seen[row.ExpirationDate]; !ok {
			dates = append(dates, row.ExpirationDate)
			seen[row.ExpirationDate] = true
		}
	}

	// Return happy
	return dates
}

//
// GetOptionsByExpirationType will return options under one expire date
//
func GetOptionsByExpirationType(expire types.Date, optionType string, options []types.OptionsChainItem) []types.OptionsChainItem {
	rt := []types.OptionsChainItem{}

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

		rt = append(rt, row)
	}

	// Return filtered subset
	return rt
}
