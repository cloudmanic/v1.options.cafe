//
// Date: 2018-10-30
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package eod

import (
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
)

//
// Get options expiration by Symbol
//
func (t *Api) GetOptionsExpirationsBySymbol(symbol string) (types.OptionsChain, error) {

	symb := strings.ToUpper(symbol)

	// New chain to return
	chain := types.OptionsChain{
		Underlying:     symb,
		ExpirationDate: types.Date{helpers.ParseDateNoError(symbol).UTC()},
		Puts:           []types.OptionsChainItem{},
		Calls:          []types.OptionsChainItem{},
	}

	// // Set the cache dir.
	// cacheDir := os.Getenv("CACHE_DIR") + "/object-store/options-eod/"

	// // Get dates from S3 store - TODO: Maybe some day just download the date we are after instead of all dates.
	// DownloadEodSymbol(symbol, false)

	// // Make sure we have this zip file
	// zipFile := cacheDir + symb + "/" + t.Day.Format("2006-01-02") + ".csv.zip"

	// // Make sure we have this file
	// if _, err := os.Stat(zipFile); os.IsNotExist(err) {
	// 	return chain, err
	// }

	// // Unzip option chain
	// f, err := files.Unzip(zipFile, "/tmp/"+symb)

	// if err != nil {
	// 	return chain, err
	// }

	// // Open CSV file
	// csvFile, err := os.Open(f[0])

	// if err != nil {
	// 	return chain, err
	// }

	// defer csvFile.Close()

	// // Read File into a Variable - https://golangcode.com/how-to-read-a-csv-file-into-a-struct
	// lines, err := csv.NewReader(csvFile).ReadAll()

	// if err != nil {
	// 	return chain, err
	// }

	// // Loop through the different lines of the CSV and Store in chain
	// for _, row := range lines {

	// 	// Get Option parts
	// 	parts, err := helpers.OptionParse(row[3])

	// 	if err != nil {
	// 		return chain, err
	// 	}

	// 	// Build Item
	// 	op := types.OptionsChainItem{
	// 		Underlying:     symb,
	// 		Symbol:         row[3],
	// 		OptionType:     parts.Type,
	// 		Description:    parts.Name,
	// 		Strike:         parts.Strike,
	// 		ExpirationDate: types.Date{parts.Expire},
	// 		Last:           helpers.StringToFloat64(row[9]),
	// 		Volume:         helpers.StringToInt(row[12]),
	// 		Bid:            helpers.StringToFloat64(row[10]),
	// 		Ask:            helpers.StringToFloat64(row[11]),
	// 		OpenInterest:   helpers.StringToInt(row[13]),
	// 		ImpliedVol:     helpers.StringToFloat64(row[11]),
	// 		Delta:          helpers.StringToFloat64(row[15]),
	// 		Gamma:          helpers.StringToFloat64(row[16]),
	// 		Theta:          helpers.StringToFloat64(row[17]),
	// 		Vega:           helpers.StringToFloat64(row[18]),
	// 	}

	// 	// Append Item
	// 	if op.OptionType == "Call" {
	// 		chain.Calls = append(chain.Calls, op)
	// 	} else if op.OptionType == "Puts" {
	// 		chain.Puts = append(chain.Puts, op)
	// 	}

	// }

	return chain, nil
}

//
// Get an options chain by expiration.
//
func (t *Api) GetOptionsChainByExpiration(symbol string, expireStr string) (types.OptionsChain, error) {

	symb := strings.ToUpper(symbol)
	expireDate := types.Date{helpers.ParseDateNoError(expireStr).UTC()}

	// New chain to return
	chain := types.OptionsChain{
		Underlying:     symb,
		ExpirationDate: expireDate,
		Puts:           []types.OptionsChainItem{},
		Calls:          []types.OptionsChainItem{},
	}

	// Get a list of all options
	options, err := t.GetOptionsBySymbol(symb)

	if err != nil {
		return chain, err
	}

	// Loop through and build chain
	for _, row := range options {

		// We only want the expire date we passed in.
		if row.ExpirationDate != expireDate {
			continue
		}

		// Append Item
		if row.OptionType == "Call" {
			chain.Calls = append(chain.Calls, row)
		} else if row.OptionType == "Put" {
			chain.Puts = append(chain.Puts, row)
		}

	}

	// Return Chain
	return chain, nil
}

/* End File */
