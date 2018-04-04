//
// Date: 2018-04-03
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-04-03
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package types

import (
	"fmt"
	"strings"
	"time"
)

// Used for date formatting with json conversion.
type Date struct {
	time.Time
}

// Used for date formatting with json conversion.
type Time struct {
	time.Time
}

// Convert JSON string to a date. Format XXXX-XX-XX
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
// Convert JSON string to a date. Make it so when we create JSON we return this format XXXX-XX-XX
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
