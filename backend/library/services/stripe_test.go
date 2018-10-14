//
// Date: 10/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nbio/st"
)

//
// Test - Apply Coupon
//
func TestApplyCoupon01(t *testing.T) {

	// Load .env file
	err := godotenv.Load("../../.env")

	if err != nil {
		panic(err)
	}

	// Create a new coupon
	couponId, err := StripeCreateNewCoupon("Unit Test Coupon 1", 55.00)
	st.Expect(t, err, nil)
	st.Expect(t, len(couponId) > 0, true)

	// Get the coupon from stripe
	coupon, err := StripeGetCoupon(couponId)

	// Verify
	st.Expect(t, err, nil)
	st.Expect(t, coupon.ID, couponId)
	st.Expect(t, coupon.Name, "Unit Test Coupon 1")
	st.Expect(t, coupon.PercentOff, 55.00)

	// Create a customer
	customerId, err := StripeAddCustomer("Jane", "Tester", "jane+unittest@options.cafe", 15)
	st.Expect(t, err, nil)

	// Add a customer to a subscription
	subId, err := StripeAddSubscription(customerId, os.Getenv("STRIPE_DEFAULT_PLAN"))
	st.Expect(t, err, nil)

	// Apply the coupon to the subscription
	err = StripeApplyCoupon(subId, couponId)
	st.Expect(t, err, nil)

	// Get the customer to verify
	customer, err := StripeGetCustomer(customerId)
	st.Expect(t, err, nil)
	st.Expect(t, customer.Subscriptions.Data[0].Discount.Coupon.ID, couponId)

	// Clean up - Delete the customer
	err = StripeDeleteCustomer(customerId)
	st.Expect(t, err, nil)

	// Clean things up by deleting the coupon
	err = StripeDeleteCoupon(couponId)
	st.Expect(t, err, nil)

}

//
// Test - CreateNewCoupon01
//
func TestCreateNewCoupon01(t *testing.T) {

	// Load .env file
	err := godotenv.Load("../../.env")

	if err != nil {
		panic(err)
	}

	// Create a new coupon
	id, err := StripeCreateNewCoupon("Unit Test Coupon 1", 55.00)
	st.Expect(t, err, nil)
	st.Expect(t, len(id) > 0, true)

	// Get the coupon from stripe
	coupon, err := StripeGetCoupon(id)

	// Verify
	st.Expect(t, err, nil)
	st.Expect(t, coupon.ID, id)
	st.Expect(t, coupon.Name, "Unit Test Coupon 1")
	st.Expect(t, coupon.PercentOff, 55.00)

	// Clean things up by deleting the coupon
	err = StripeDeleteCoupon(id)
	st.Expect(t, err, nil)

}

/* End File */
