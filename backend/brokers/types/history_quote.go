//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package types

import (
	"fmt"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

type Date struct {
	time.Time
}

type HistoryQuote struct {
	Date   Date    `json:"date"`
	Time   Time    `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
}

//
// Convert json string to a date.
//
func (t *Date) UnmarshalJSON(b []byte) error {

	// Remove quotes
	str := strings.Replace(string(b), "\"", "", -1)

	// Parse string
	tt, _ := time.Parse("2006-01-02", str)

	// Return UTC
	*t = Date{tt.UTC()}
	return nil
}

//
// Convert json string to a date.
//
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))), nil
}

//
// Convert json string to a time.
//
func (t *Time) UnmarshalJSON(b []byte) error {
	// Remove quotes
	str := strings.Replace(string(b), "\"", "", -1)

	// Add Timezone (hard coded EST)
	str = str + "-0500 EST"

	// Parse string
	tt, _ := time.Parse("2006-01-02T15:04:05-0700 MST", str)

	// Return UTC
	*t = Time{tt.UTC()}
	return nil
}

/* End File */
