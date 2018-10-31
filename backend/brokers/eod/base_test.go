//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"testing"

	"github.com/nbio/st"
)

//
// Test - GetTradeDatesBySymbols01
//
func TestGetTradeDatesBySymbols01(t *testing.T) {

	// Get dates from S3 store
	dates, err := GetTradeDatesBySymbols("spy")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, (len(dates) >= 3477), true)

}

//
// Test - GetTradeDateKeysBySymbol01
//
func TestGetTradeDateKeysBySymbol01(t *testing.T) {

	// Get dates from S3 store
	keys, err := GetTradeDateKeysBySymbol("spy")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, (len(keys) >= 3477), true)

}

//
// Test - DownloadEodSymbol01
//
func TestDownloadEodSymbol01(t *testing.T) {

	// Get dates from S3 store
	files := DownloadEodSymbol("spy", false)

	// Test result
	st.Expect(t, (len(files) >= 3477), true)

}

/* End File */
