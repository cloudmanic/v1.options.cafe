//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// Notes: the reason this is in the services package is we don't
// want to conflict with the stripe name space
//

package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balance"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/coupon"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

//
// Get one transaction balance.
//
func StripeGetBalanceTransaction(id string) (*stripe.BalanceTransaction, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction")
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get the transaction
	bt, err := balance.GetBalanceTransaction(id, nil)

	if err != nil {
		BetterError(errors.New("StripeGetBalanceTransaction : Unable to get a transaction balance. " + id + " (" + err.Error() + ")"))
		return nil, err
	}

	// Return happy
	return bt, nil

}

//
// Apply a coupon to a Subscription
//
func StripeApplyCoupon(subId string, couponId string) error {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeDeleteCoupon")
		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteNewCoupon")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.SubscriptionParams{
		Coupon: stripe.String(couponId),
	}

	// Send request to stripe
	_, err := sub.Update(subId, params)

	if err != nil {
		BetterError(errors.New("StripeApplyCoupon : Unable to apply a coupon. " + couponId + " (" + err.Error() + ")"))
		return err
	}

	// Return happy
	return nil

}

//
// Get Coupon.
//
func StripeGetCoupon(id string) (*stripe.Coupon, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeAddCustomer")
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetCoupon")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object and Delete
	c, err := coupon.Get(id, nil)

	if err != nil {
		BetterError(errors.New("StripeGetCoupon : Unable to get a coupon. " + id + " (" + err.Error() + ")"))
		return nil, err
	}

	// Return happy
	return c, nil

}

//
// Delete a coupon.
//
func StripeDeleteCoupon(id string) error {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeDeleteCoupon")
		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteNewCoupon")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Delete coupon.
	_, err := coupon.Del(id, nil)

	if err != nil {
		BetterError(errors.New("StripeDeleteCoupon : Unable to create a new coupon. (" + err.Error() + ")"))
		return err
	}

	// Return the new coupon Id
	return nil
}

//
// Create a new coupon.
//
func StripeCreateNewCoupon(name string, percentOff float64) (string, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in CreateNewCoupon")
		return "", errors.New("No STRIPE_SECRET_KEY found in StripeCreateNewCoupon")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the coupon object
	params := &stripe.CouponParams{
		PercentOff: stripe.Float64(percentOff),
		Duration:   stripe.String(string(stripe.CouponDurationForever)),
		Name:       stripe.String(name),
	}

	// Create new coupon.
	couponObj, err := coupon.New(params)

	if err != nil {
		BetterError(errors.New("StripeCreateNewCoupon : Unable to create a new coupon. (" + err.Error() + ")"))
		return "", err
	}

	// Return the new coupon Id
	return couponObj.ID, nil
}

//
// Add a new customer.
//
func StripeAddCustomer(first string, last string, email string, accountId int) (string, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeAddCustomer")
		return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddCustomer")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object
	customerParams := &stripe.CustomerParams{Email: &email}
	customerParams.AddMetadata("FirstName", first)
	customerParams.AddMetadata("LastName", last)
	customerParams.AddMetadata("AccountId", strconv.Itoa(accountId))

	// Create new customer.
	customer, err := customer.New(customerParams)

	if err != nil {
		BetterError(errors.New("StripeAddCustomer : Unable to create a new customer. (" + err.Error() + ")"))
		return "", err
	}

	// Return the new customer Id
	return customer.ID, nil

}

//
// Delete a new customer.
//
func StripeDeleteCustomer(custToken string) error {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeAddCustomer")
		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCustomer")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object and Delete
	params := &stripe.CustomerParams{}
	_, err := customer.Del(custToken, params)

	if err != nil {
		BetterError(errors.New("StripeDeleteCustomer : Unable to delete a customer. " + custToken + " (" + err.Error() + ")"))
		return err
	}

	// Return happy
	return nil

}

//
// Get customer.
//
func StripeGetCustomer(custToken string) (*stripe.Customer, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeAddCustomer")
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCustomer")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object and Delete
	cust, err := customer.Get(custToken, nil)

	if err != nil {
		BetterError(errors.New("StripeDeleteCustomer : Unable to get a customer. " + custToken + " (" + err.Error() + ")"))
		return nil, err
	}

	// Return happy
	return cust, nil

}

//
// Add a customer subscription.
//
func StripeAddSubscription(custId string, plan string, coupon string) (string, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in StripeAddSubscription")
		return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddSubscription")
	}

	// Setup Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(custId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(plan),
			},
		},
	}

	// Do we have a coupon code?
	if len(coupon) > 0 {
		subParams.Coupon = stripe.String(coupon)
		Info("Coupon code passed with subscribe token: " + coupon + " - " + custId)
	}

	// Create new subscription.
	subscription, err := sub.New(subParams)

	if err != nil {
		BetterError(errors.New("StripeAddSubscription : Unable to create a new subscription. (" + err.Error() + ")"))
		return "", err
	}

	// Return the new subscription Id
	return subscription.ID, nil

}

//
// Add a new credit card.
//
func StripeAddCreditCardByToken(custId string, token string) (string, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in AddCreditCardByToken")
		return "", errors.New("No STRIPE_SECRET_KEY found in AddCreditCardByToken")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CardParams{
		Customer: stripe.String(custId),
		Token:    stripe.String(token),
	}

	// call to stripe to add card
	c, err := card.New(params)

	if err != nil {
		BetterError(errors.New("AddCreditCardByToken : Unable to add card. (" + err.Error() + ")"))
		return "", err
	}

	// Return the new card Id
	return c.ID, nil

}

//
// List all credit cards on file.
//
func StripeListAllCreditCards(custId string) ([]string, error) {

	cardList := []string{}

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in ListAllCreditCards")
		return nil, errors.New("No STRIPE_SECRET_KEY found in ListAllCreditCards")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CardListParams{
		Customer: stripe.String(custId),
	}

	params.Filters.AddFilter("limit", "", "100")

	list := card.List(params)

	for list.Next() {
		c := list.Card()
		cardList = append(cardList, c.ID)
	}

	// Return the card list
	return cardList, nil

}

//
// Delete credit cards on file.
//
func StripeDeleteCreditCard(custId string, cardId string) error {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		Critical("No STRIPE_SECRET_KEY found in DeleteCreditCard")
		return errors.New("No STRIPE_SECRET_KEY found in DeleteCreditCard")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CardParams{
		Customer: stripe.String(custId),
	}

	// Delete card at stripe
	_, err := card.Del(cardId, params)

	if err != nil {
		BetterError(errors.New("DeleteCreditCard : Unable to delete card. (" + err.Error() + ")"))
		return err
	}

	// Return happy
	return nil

}

/* End File */
