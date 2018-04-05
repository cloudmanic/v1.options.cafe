//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

type OptionParts struct {
	Name   string
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
func OptionParse(optionSymb string) (OptionParts, error) {

	// Length must be at least 16 chars
	if len(optionSymb) < 16 {
		return OptionParts{}, errors.New("Options symbol string must at least 16 chars. - " + optionSymb)
	}

	// Parse the options string
	// https://regex101.com/r/jEGDzO/1
	var re = regexp.MustCompile(`(\D{1,6})(\d{2})(\d{2})(\d{2})(C|P)(\d{5})(\d{3})`)

	// Get the matches
	matches := re.FindAllStringSubmatch(optionSymb, -1)

	// Convert to ints
	year, _ := strconv.Atoi(matches[0][2])
	month, _ := strconv.Atoi(matches[0][3])
	day, _ := strconv.Atoi(matches[0][4])
	year = year + 2000 // TODO: in year 3000 this will be a bug :)

	// Build date
	date, err := dateparse.ParseLocal(matches[0][3] + "/" + matches[0][4] + "/" + matches[0][2])

	if err != nil {
		return OptionParts{}, err
	}

	// Get type string.
	var typeStr = "Call"

	if matches[0][5] == "P" {
		typeStr = "Put"
	}

	// Build strike
	baseNum, _ := strconv.Atoi(matches[0][6])
	decimalNum, _ := strconv.Atoi(matches[0][7])
	strike := float64(baseNum) + (float64(decimalNum) * .001)

	// Build the name string of the option (ie. "SPY Mar 16 2018 $250.00 Put")
	name := matches[0][1] + " " + date.Format("Jan 2, 2006") + " $" + fmt.Sprintf("%.2f", strike) + " " + typeStr

	// Returned parsed object
	return OptionParts{
		Name:   name,
		Option: optionSymb,
		Symbol: matches[0][1],
		Year:   uint(year),
		Month:  uint(month),
		Day:    uint(day),
		Expire: date.UTC(),
		Type:   typeStr,
		Strike: strike,
	}, nil
}

/* End File */
