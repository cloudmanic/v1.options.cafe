//
// Date: 2/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cache

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test - Set 01
//
func TestSet01(t *testing.T) {

	// Store something in cache
	Set("oc-testing-1", "Options Cafe is DaBomb.com")

	// Get stored value.
	result := ""
	found, err := Get("oc-testing-1", &result)

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, found, true)
	st.Expect(t, result, "Options Cafe is DaBomb.com")
}

//
// Test - Set 02
//
func TestSet02(t *testing.T) {

	// Get a value we know we do not have
	result := ""
	found, _ := Get("oc-testing-not-found", &result)

	// Verify the data was return as expected
	st.Expect(t, found, false)
	st.Expect(t, result, "")
}

//
// Test - Set 03
//
func TestSet03(t *testing.T) {

	// Create an order model.
	order := models.Order{
		Type:           "Type Here",
		SymbolId:       uint(1),
		OptionSymbolId: uint(1),
		Side:           "Side Here",
		Qty:            10,
		Status:         "Open",
		Legs: []models.OrderLeg{
			{Type: "Type OrderLeg 1", Side: "Side OrderLeg 1"},
			{Type: "Type OrderLeg 2", Side: "Side OrderLeg 2"},
		},
	}

	// Store the struct in the cache
	Set("oc-testing-2", order)

	// Get a value we know we do not have
	result := models.Order{}
	found, _ := Get("oc-testing-2", &result)

	// Verify the data was return as expected
	st.Expect(t, found, true)
	st.Expect(t, result.Type, "Type Here")
	st.Expect(t, result.Legs[0].Type, "Type OrderLeg 1")
	st.Expect(t, result.Legs[1].Type, "Type OrderLeg 2")
}

/* End File */
