//
// Date: 2/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"

	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test - Get Orders By User Class Status Reviewed
//
func TestGetOrdersByUserClassStatusReviewed01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Setup know test parms
	var userId uint = 1
	var class string = "multileg"
	var status string = "filled"
	var reviewed string = "No"

	// Make query.
	orders, err := db.GetOrdersByUserClassStatusReviewed(userId, class, status, reviewed)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, len(orders), 15)
	st.Expect(t, orders[4].Id, uint(9))
	st.Expect(t, len(orders[4].Legs), 2)
	st.Expect(t, orders[4].Legs[0].Id, uint(13))
	st.Expect(t, orders[4].Legs[1].Id, uint(14))
}

/* End File */
