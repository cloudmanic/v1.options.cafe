//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
  "os"
  "errors"
  "strconv"  
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/sub"
  "github.com/stripe/stripe-go/customer" 
)

//
// Add a new customer.
//
func AddCustomer(first string, last string, email string, accountId int) (string, error) {
  
  // Make sure we have a STRIPE_SECRET_KEY
  if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
    MajorLog("No STRIPE_SECRET_KEY found in StripeAddCustomer")
    return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddCustomer")
  }  
  
  stripe.Key = os.Getenv("STRIPE_SECRET_KEY") 
  
  // Setup the customer object
	customerParams := &stripe.CustomerParams{ Email: email }
	customerParams.AddMeta("FirstName", first)  
	customerParams.AddMeta("LastName", last) 
	customerParams.AddMeta("AccountId", strconv.Itoa(accountId))
  
  // Create new customer.
	customer, err := customer.New(customerParams)

	if err != nil {
  	Error(err, "StripeAddCustomer : Unable to create a new customer.")
    return "", err
	}
	
	// Return the new customer Id
	return customer.ID, nil
  
}

//
// Add a customer subscription.
//
func AddSubscription(custId string, plan string) (string, error) {
  
  // Make sure we have a STRIPE_SECRET_KEY
  if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
    MajorLog("No STRIPE_SECRET_KEY found in StripeAddSubscription")
    return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddSubscription")
  }   
  
  stripe.Key = os.Getenv("STRIPE_SECRET_KEY") 
  
  // Setup the customer object
	subParams := &stripe.SubParams{ Customer: custId, Plan: plan }
  
  // Create new subscription.
	subscription, err := sub.New(subParams)

	if err != nil {
  	Error(err, "StripeAddSubscription : Unable to create a new subscription.")
    return "", err
	}
	
	// Return the new subscription Id
	return subscription.ID, nil  
  
}

/* End File */