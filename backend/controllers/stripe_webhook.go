//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
  "os"
  "fmt"
  "net/http"
  "io/ioutil"
  "app.options.cafe/backend/models"
  "github.com/stripe/stripe-go/webhook"
  "app.options.cafe/backend/library/services"
)

//
// Manage webhooks from strip.
//
func DoStripeWebhook(w http.ResponseWriter, r *http.Request) {
  
  // Get body of the request.
  body, err := ioutil.ReadAll(r.Body)
  
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  
  // Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
  event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_SIGNING_SECRET"))
  
  if err != nil {
    w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
    fmt.Fprintf(w, "%v", err)
    return
  }
    
  defer r.Body.Close()
  
  // Log the event.
  services.Log("Stripe Webhook Received : " + event.Type + " - " + event.ID)
  
  // Start the database
  DB.Start()
  defer DB.Connection.Close()  
  
  // Figure out what user this event is for.
  // event.GetObjValue("customer")
  user, err := DB.GetUserByStripeCustomer("cus_BKQjz6z9hMqzYY")
  
  if err != nil {
    services.MajorLog("Stripe Webhook Unknown user found in event : " + event.Type + " - " + event.ID)
    w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
    fmt.Fprintf(w, "%v", err)
    return 
  }
  
  // Figure out how we handle this event.
  switch event.Type {
    
    // customer.subscription.deleted
    case "customer.subscription.deleted":
      StripeEventSubscriptionDeleted(user)
      
  }
  
  // Tell stripe all went well.
  fmt.Fprintf(w, "Done")
   
}

// -------------------- Stripe Events --------------------------- //

//
// Stripe event : customer.subscription.deleted
//
func StripeEventSubscriptionDeleted(user models.User) {
  
  // Remove the subscription from the database.
  user.StripeSubscription = ""
  DB.Connection.Save(&user)
  
  // Log event.
  services.Log("Stripe Subscription Deleted: " + user.Email); 
  
}

/* End File */