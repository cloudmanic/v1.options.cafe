//
// Date: 2/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"testing"

	"github.com/nbio/st"
)

//
// Test - Option parse 01
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
	st.Expect(t, parts.Expire.Format("2006-01-02"), "2018-02-23")
	st.Expect(t, parts.Type, "Call")
	st.Expect(t, parts.Strike, float64(55))
}

//
// Test - Option parse 02
//
func TestOptionParse02(t *testing.T) {

	// Test the parser
	parts, err := OptionParse("SPY180803C00295000")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, parts.Name, "SPY Aug 3, 2018 $295.00 Call")
	st.Expect(t, parts.Option, "SPY180803C00295000")
	st.Expect(t, parts.Symbol, "SPY")
	st.Expect(t, parts.Year, uint(2018))
	st.Expect(t, parts.Month, uint(8))
	st.Expect(t, parts.Day, uint(3))
	st.Expect(t, parts.Expire.Format("2006-01-02"), "2018-08-03")
	st.Expect(t, parts.Type, "Call")
	st.Expect(t, parts.Strike, float64(295))
}

/* End File */
