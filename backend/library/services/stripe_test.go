//
// Date: 10/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"go/build"
	"os"
	"testing"

	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test getting a GetChargesByCustomer
//
func TestGetChargesByCustomer01(t *testing.T) {

	// // We comment this out because "cus_Djqq5q9mnW0lm6" might not be reliable.
	// // We could do better testing like create a charge and test

	// // Load .env file
	// err := godotenv.Load("../../.env")

	// if err != nil {
	// 	panic(err)
	// }

	// // Get transaction balance
	// ch, err := StripeGetChargesByCustomer("cus_Djqq5q9mnW0lm6")
	// st.Expect(t, err, nil)

	// //spew.Dump(ch)

	// for _, row := range ch {

	// 	if row.Invoice != nil {
	// 		inv, _ := StripeGetInvoice(row.Invoice.ID)
	// 		spew.Dump(inv)
	// 	}

	// }
}

//
// Test getting a BalanceTransaction
//
func TestGetBalanceTransaction01(t *testing.T) {

	// We comment this out because "txn_1DLGzaKah1jS67deeMUpKcBG" might not be reliable.
	// We could do better testing like create a charge and then grab the Balance transaction

	// // Load .env file
	// err := godotenv.Load("../../.env")

	// if err != nil {
	// 	panic(err)
	// }

	// // Get transaction balance
	// tb, err := StripeGetBalanceTransaction("txn_1DLGzaKah1jS67deeMUpKcBG")
	// st.Expect(t, err, nil)
}

//
// Test - Apply Coupon
//
func TestApplyCoupon01(t *testing.T) {

	// Load .env file
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")

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

	// Add Payment Source
	_, err = StripeAddCreditCardByToken(customerId, "tok_visa")
	st.Expect(t, err, nil)

	// Add a customer to a subscription
	subId, err := StripeAddSubscription(customerId, os.Getenv("STRIPE_MONTHLY_PLAN"), "")
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
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")

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
