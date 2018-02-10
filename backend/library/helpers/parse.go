//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"regexp"
	"strconv"
	"time"
)

type OptionParts struct {
	Option string
	Symbol string
	Year   uint
	Month  uint
	Day    uint
	Expire time.Time
	Type   string
	Strike float64
}

//
// Parse the an options symbol string.
// https://www.investopedia.com/articles/optioninvestor/10/options-symbol-rules.asp
//
func OptionParse(optionSymb string) OptionParts {

	// Parse the options string
	// https://regex101.com/r/jEGDzO/1
	var re = regexp.MustCompile(`(\D{1,6})(\d{2})(\d{2})(\d{2})(C|P)(\d{5})(\d{3})`)

	// Get the matches
	matches := re.FindAllStringSubmatch(optionSymb, -1)

	// Convert to ints
	year, _ := strconv.Atoi(matches[0][2])
	month, _ := strconv.Atoi(matches[0][3])
	day, _ := strconv.Atoi(matches[0][4])

	// Build date
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Get type string.
	var typeStr = "Call"

	if matches[0][5] == "P" {
		typeStr = "Put"
	}

	// Build strike
	baseNum, _ := strconv.Atoi(matches[0][6])
	decimalNum, _ := strconv.Atoi(matches[0][7])
	strike := float64(baseNum) + (float64(decimalNum) * .001)

	// Returned parsed object
	return OptionParts{
		Option: optionSymb,
		Symbol: matches[0][1],
		Year:   uint(year),
		Month:  uint(month),
		Day:    uint(day),
		Expire: date,
		Type:   typeStr,
		Strike: strike,
	}
}

/* End File */
