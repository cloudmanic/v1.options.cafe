//
// Date: 2/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"testing"
	"time"

	"github.com/nbio/st"
)

//
// Test - Option parse
//
func TestOptionParse01(t *testing.T) {

	// Test the parser
	parts, err := OptionParse("VXX180223C00055000")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, parts.Name, "VXX Feb 23, 2018 $55.00 Call")
	st.Expect(t, parts.Option, "VXX180223C00055000")
	st.Expect(t, parts.Symbol, "VXX")
	st.Expect(t, parts.Year, uint(2018))
	st.Expect(t, parts.Month, uint(2))
	st.Expect(t, parts.Day, uint(23))
	st.Expect(t, parts.Expire, time.Date(2018, 02, 23, 00, 00, 00, 0000, time.UTC))
	st.Expect(t, parts.Type, "Call")
	st.Expect(t, parts.Strike, float64(55))
}

/* End File */
