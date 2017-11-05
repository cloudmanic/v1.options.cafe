//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"app.options.cafe/backend/emails"
	"app.options.cafe/backend/library/email"
	"app.options.cafe/backend/library/services"
	"app.options.cafe/backend/models"
	"github.com/stripe/stripe-go/webhook"
)

//
// Manage webhooks from strip.
//
func (t *Controller) DoStripeWebhook(w http.ResponseWriter, r *http.Request) {

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

	// If there is no customer value there is nothing we need to do.
	// As of now all the events we care about have a customer attached to them.
	if len(event.GetObjValue("customer")) == 0 {
		services.Log("Stripe no customer data found (this is expected).")
		fmt.Fprintf(w, "Done")
		return
	}

	// Figure out what user this event is for.
	user, err := t.DB.GetUserByStripeCustomer(event.GetObjValue("customer"))

	if err != nil {
		services.MajorLog("Stripe Webhook Unknown user found in event : " + event.Type + " - " + event.ID)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error
		fmt.Fprintf(w, "%v", err)
		return
	}

	// Figure out how we handle this event.
	switch event.Type {

	// customer.subscription.deleted
	case "customer.subscription.deleted":
		t.StripeEventSubscriptionDeleted(user)

	}

	// Tell stripe all went well.
	fmt.Fprintf(w, "Done")

}

// -------------------- Stripe Events --------------------------- //

//
// Stripe event : customer.subscription.deleted
//
func (t *Controller) StripeEventSubscriptionDeleted(user models.User) {

	// Remove the subscription from the database.
	user.StripeSubscription = ""
	t.DB.UpdateUser(&user)

	// Log event.
	services.Log("Stripe Subscription Deleted: " + user.Email)

	// Send email telling the user this happened.
	var url = "https://app.options.cafe"

	go email.Send(
		user.Email,
		"Options Cafe : Subscription Canceled",
		emails.GetSubscriptionCanceledHtml(user.FirstName, url),
		emails.GetSubscriptionCanceledText(user.FirstName, url))

	// Tell slack about this.
	go services.SlackNotify("#events", "Options Cafe Subscription Canceled : "+user.Email)

}

/* End File */
