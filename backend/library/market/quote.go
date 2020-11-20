//
// Date: 2018-12-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package market

import (
	"errors"
	"time"

	"app.options.cafe/brokers"
	"app.options.cafe/models"
)

//
// GetUnderlayingQuoteByDate - Return a underlying stock quote based on date in time.
//
func GetUnderlayingQuoteByDate(db models.Datastore, userId uint, symbol string, today time.Time) (float64, error) {
	// Special case for some legacy symbols (Thanks Tradier)
	if symbol == "SPXW" {
		symbol = "SPX"
	}

	if symbol == "RUTW" {
		symbol = "RUT"
	}

	// See if we already have this quote in our DB
	q := models.HistoricalQuote{}
	db.New().Where("short_name = ? AND date = ?", symbol, today.Format("2006-01-02")).First(&q)

	// If we have the quote just returh this and do not call broker API
	if q.Id > 0 {
		return q.Price, nil
	}

	// Get a broker to use to get this data
	broker, err := brokers.GetPrimaryTradierConnection(db, userId)

	if err != nil {
		return 0.00, err
	}

	// Set start date
	start, err := time.Parse("2006-01-02 15:04", today.Format("2006-01-02")+" 09:30")

	if err != nil {
		return 0.00, err
	}

	// Set end date
	end, _ := time.Parse("2006-01-02 15:04", today.Format("2006-01-02")+" 16:00")

	if err != nil {
		return 0.00, err
	}

	// Make call to broker.
	result, err := broker.GetHistoricalQuotes(symbol, start, end, "daily")

	if err != nil {
		return 0.00, err
	}

	if len(result) < 1 {
		return 0.00, errors.New("GetUnderlayingQuoteByDate: No result found.")
	}

	// Setup historical quote to store
	h := models.HistoricalQuote{
		ShortName: symbol,
		Date:      models.Date{today},
		Price:     result[0].Close,
	}

	// Store in DB
	db.New().Where("short_name = ? AND date = ?", symbol, today.Format("2006-01-02")).FirstOrCreate(&h)

	// Return Happy
	return result[0].Close, nil
}

/* End File */
