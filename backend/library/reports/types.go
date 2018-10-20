//
// Date: 2/28/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"fmt"
	"strings"
	"time"
)

// Used for date formatting with json conversion.
type Date struct {
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

/* End File */
