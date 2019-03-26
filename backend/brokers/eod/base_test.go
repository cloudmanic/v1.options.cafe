//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-07
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"os"
	"testing"

	"github.com/nbio/st"

	"github.com/cloudmanic/app.options.cafe/backend/library/files"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Test - GetGetOptionsBySymbol01
//
func TestGetOptionsBySymbol01(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping GetGetOptionsBySymbol01 test since it requires a broker token and --short was requested")
	}

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create broker object
	o := Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// Get options
	options, underlyingLast, err := o.GetOptionsBySymbol("spy")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, underlyingLast, 276.39)
	st.Expect(t, len(options), 6262)
}

//
// Test - GetTradeDatesBySymbols01
//
func TestGetTradeDatesBySymbols01(t *testing.T) {

	if testing.Short() {
		t.Skipf("Skipping TestGetTradeDatesBySymbols01 test since it requires a broker token and --short was requested")
	}

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

	if testing.Short() {
		t.Skipf("Skipping TestGetTradeDateKeysBySymbol01 test since it requires a broker token and --short was requested")
	}

	// Get dates from S3 store
	keys, err := GetTradeDateKeysBySymbol("spy")

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, (len(keys) >= 3477), true)

}

//
// Test - downloadEodSymbol01
//
func TestDownloadEodSymbol01(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping TestDownloadEodSymbol01 test since it requires a broker token and --short was requested")
	}

	// Download File
	dFile := "options-eod/SPY/2018-10-18.csv.zip"
	lFile := os.Getenv("CACHE_DIR") + "/object-store/" + dFile

	// Delete any left overs of this file so we can test the download.
	os.Remove(lFile)

	// Get dates from S3 store
	file, err := downloadEodSymbol("spy", helpers.ParseDateNoError("2018-10-18").UTC())

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, file, lFile)
	st.Expect(t, "f39b09fc32df9640ee1f4aa21f74ae93", files.Md5(lFile))
}

//
// Test - TestUnzipSymbolCSV01
//
func TestUnzipSymbolCSV01(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping TestUnzipSymbolCSV01 test since it requires a broker token and --short was requested")
	}

	// Download File
	dFile := "options-eod/SPY/2018-10-18.csv.zip"
	lFile := os.Getenv("CACHE_DIR") + "/object-store/" + dFile

	// Get dates from S3 store
	file, err := downloadEodSymbol("spy", helpers.ParseDateNoError("2018-10-18").UTC())

	// Just double check.
	st.Expect(t, err, nil)
	st.Expect(t, file, lFile)
	st.Expect(t, "f39b09fc32df9640ee1f4aa21f74ae93", files.Md5(lFile))

	// Unzip file.
	options, last, err := unzipSymbolCSV("spy", file)

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, 276.39, last)
	st.Expect(t, len(options), 6262)
	st.Expect(t, options[55].Underlying, "SPY")
	st.Expect(t, options[55].Last, 0.01)
	st.Expect(t, options[55].Symbol, "SPY181019P00199000")
}

/* End File */
