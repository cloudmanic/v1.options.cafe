//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"testing"
	"time"

	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - GetHistoricalQuotes
//
func TestGetHistoricalQuotes01(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/history").
		Reply(200).
		BodyString(`{"history":{"day":[{"date":"2018-01-02","open":267.84,"high":268.81,"low":267.4,"close":268.77,"volume":86631393},{"date":"2018-01-03","open":268.96,"high":270.64,"low":268.96,"close":270.47,"volume":90054249},{"date":"2018-01-04","open":271.2,"high":272.16,"low":270.5447,"close":271.61,"volume":80624896},{"date":"2018-01-05","open":272.51,"high":273.56,"low":271.95,"close":273.42,"volume":83522436},{"date":"2018-01-08","open":273.31,"high":274.1,"low":272.98,"close":273.92,"volume":57310300},{"date":"2018-01-09","open":274.4,"high":275.25,"low":274.081,"close":274.54,"volume":57245653},{"date":"2018-01-10","open":273.68,"high":274.42,"low":272.92,"close":274.12,"volume":69543174},{"date":"2018-01-11","open":274.75,"high":276.12,"low":274.56,"close":276.12,"volume":62357781},{"date":"2018-01-12","open":276.42,"high":278.11,"low":276.0819,"close":277.92,"volume":90807258},{"date":"2018-01-16","open":279.35,"high":280.09,"low":276.18,"close":276.97,"volume":106532827},{"date":"2018-01-17","open":278.03,"high":280.05,"low":276.97,"close":279.61,"volume":113242394},{"date":"2018-01-18","open":279.48,"high":279.96,"low":278.58,"close":279.14,"volume":100723726},{"date":"2018-01-19","open":279.8,"high":280.41,"low":279.14,"close":280.41,"volume":140912870},{"date":"2018-01-22","open":280.17,"high":282.69,"low":280.11,"close":282.69,"volume":91302336},{"date":"2018-01-23","open":282.74,"high":283.62,"low":282.37,"close":283.29,"volume":97073068},{"date":"2018-01-24","open":284.02,"high":284.7,"low":281.84,"close":283.18,"volume":134821100},{"date":"2018-01-25","open":284.16,"high":284.27,"low":282.405,"close":283.3,"volume":84573237},{"date":"2018-01-26","open":284.25,"high":286.6285,"low":283.96,"close":286.58,"volume":107729783},{"date":"2018-01-29","open":285.93,"high":286.43,"low":284.5,"close":284.68,"volume":90100642},{"date":"2018-01-30","open":282.6,"high":284.736,"low":281.22,"close":281.76,"volume":131751380},{"date":"2018-01-31","open":282.73,"high":283.3,"low":280.68,"close":281.9,"volume":118923457},{"date":"2018-02-01","open":281.07,"high":283.06,"low":280.68,"close":281.58,"volume":90083692},{"date":"2018-02-02","open":280.08,"high":280.23,"low":275.41,"close":275.45,"volume":173164270},{"date":"2018-02-05","open":273.45,"high":275.85,"low":263.31,"close":263.93,"volume":294568480},{"date":"2018-02-06","open":259.94,"high":269.7,"low":258.7,"close":269.13,"volume":355018200},{"date":"2018-02-07","open":268.5,"high":272.36,"low":267.58,"close":267.67,"volume":167241140},{"date":"2018-02-08","open":268.01,"high":268.17,"low":257.59,"close":257.63,"volume":246397060},{"date":"2018-02-09","open":260.8,"high":263.61,"low":252.92,"close":261.5,"volume":283490830},{"date":"2018-02-12","open":263.83,"high":267.01,"low":261.6644,"close":265.34,"volume":143725730},{"date":"2018-02-13","open":263.97,"high":266.62,"low":263.31,"close":266.0,"volume":81210909},{"date":"2018-02-14","open":264.31,"high":270.0,"low":264.3,"close":269.59,"volume":120683003},{"date":"2018-02-15","open":271.57,"high":273.04,"low":268.77,"close":273.03,"volume":111174904},{"date":"2018-02-16","open":272.32,"high":275.32,"low":272.27,"close":273.11,"volume":160417710},{"date":"2018-02-20","open":272.03,"high":273.67,"low":270.5,"close":271.4,"volume":86285058},{"date":"2018-02-21","open":271.9,"high":274.72,"low":269.94,"close":270.05,"volume":98714168},{"date":"2018-02-22","open":271.1,"high":273.05,"low":269.64,"close":270.4,"volume":110433492},{"date":"2018-02-23","open":271.79,"high":274.71,"low":271.25,"close":274.71,"volume":92753578}]}}`)

	// Create new tradier instance
	tradier := &Api{}

	// Send API call.
	start := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2018, 02, 23, 0, 0, 0, 0, time.UTC)
	quotes, err := tradier.GetHistoricalQuotes("SPY", start, end, "daily")

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(quotes), 37)
	st.Expect(t, quotes[0].Date, "2018-01-02")
	st.Expect(t, quotes[0].Open, 267.84)
	st.Expect(t, quotes[0].High, 268.81)
	st.Expect(t, quotes[0].Low, 267.4)
	st.Expect(t, quotes[0].Close, 268.77)
	st.Expect(t, quotes[0].Volume, 86631393)
	st.Expect(t, quotes[36].Date, "2018-02-23")
	st.Expect(t, quotes[36].Open, 271.79)
	st.Expect(t, quotes[36].High, 274.71)
	st.Expect(t, quotes[36].Low, 271.25)
	st.Expect(t, quotes[36].Close, 274.71)
	st.Expect(t, quotes[36].Volume, 92753578)

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

//
// Test - GetHistoricalQuotes
//
func TestGetHistoricalQuotes02(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/history").
		Reply(200).
		BodyString(`{"history":{"day":[{"date":"2017-01-01","open":135.1,"high":139.37,"low":133.05,"close":138.47,"volume":79912670},{"date":"2017-02-01","open":137.66,"high":146.33,"low":136.33,"close":144.91,"volume":75930406},{"date":"2017-03-01","open":146.72,"high":150.15,"low":145.83,"close":146.95,"volume":89749218},{"date":"2017-04-01","open":146.94,"high":156.27,"low":145.76,"close":156.1,"volume":75717119},{"date":"2017-05-01","open":156.22,"high":160.86,"low":153.29,"close":153.88,"volume":86552189},{"date":"2017-06-01","open":153.52,"high":159.2242,"low":150.75,"close":153.4,"volume":100510570},{"date":"2017-07-01","open":154.39,"high":154.79,"low":144.25,"close":148.08,"volume":97247580},{"date":"2017-08-01","open":150.24,"high":156.05,"low":146.89,"close":150.0,"volume":108301722},{"date":"2017-09-01","open":150.26,"high":163.61,"low":149.76,"close":163.56,"volume":110304258},{"date":"2017-10-01","open":164.2,"high":167.94,"low":161.5111,"close":165.31,"volume":68423735},{"date":"2017-11-01","open":166.42,"high":180.67,"low":160.53,"close":179.82,"volume":103228255},{"date":"2017-12-01","open":180.32,"high":191.49,"low":176.7,"close":189.53,"volume":103713592},{"date":"2018-01-01","open":190.21,"high":207.605,"low":187.819,"close":201.81,"volume":87523559},{"date":"2018-02-01","open":199.34,"high":202.25,"low":175.42,"close":188.35,"volume":111317234}]}}`)

	// Create new tradier instance
	tradier := &Api{}

	// Send API call.
	start := time.Date(2017, 01, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2018, 02, 23, 0, 0, 0, 0, time.UTC)
	quotes, err := tradier.GetHistoricalQuotes("HD", start, end, "monthly")

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(quotes), 14)
	st.Expect(t, quotes[0].Date, "2017-01-01")
	st.Expect(t, quotes[0].Open, 135.1)
	st.Expect(t, quotes[0].High, 139.37)
	st.Expect(t, quotes[0].Low, 133.05)
	st.Expect(t, quotes[0].Close, 138.47)
	st.Expect(t, quotes[0].Volume, 79912670)
	st.Expect(t, quotes[13].Date, "2018-02-01")
	st.Expect(t, quotes[13].Open, 199.34)
	st.Expect(t, quotes[13].High, 202.25)
	st.Expect(t, quotes[13].Low, 175.42)
	st.Expect(t, quotes[13].Close, 188.35)
	st.Expect(t, quotes[13].Volume, 111317234)

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

/* End File */
