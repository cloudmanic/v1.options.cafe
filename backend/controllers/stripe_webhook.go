//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"os"

	"github.com/app.options.cafe/backend/emails"
	"github.com/app.options.cafe/backend/library/email"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/webhook"
)

//
// Manage webhooks from strip.
//
func (t *Controller) DoStripeWebhook(c *gin.Context) {

	// Get body of the request.
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
	event, err := webhook.ConstructEvent(body, c.Request.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_SIGNING_SECRET"))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	defer c.Request.Body.Close()

	// Log the event.
	services.Info("Stripe Webhook Received : " + event.Type + " - " + event.ID)

	// If there is no customer value there is nothing we need to do.
	// As of now all the events we care about have a customer attached to them.
	if len(event.GetObjValue("customer")) == 0 {
		services.Info("Stripe no customer data found (this is expected).")
		c.JSON(200, gin.H{"status": "success"})
		return
	}

	// Figure out what user this event is for.
	user, err := t.DB.GetUserByStripeCustomer(event.GetObjValue("customer"))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Figure out how we handle this event.
	switch event.Type {

	// customer.subscription.deleted
	case "customer.subscription.deleted":
		t.StripeEventSubscriptionDeleted(user)

	}

	// Tell stripe all went well.
	c.JSON(200, gin.H{"status": "success"})
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
	services.Info("Stripe Subscription Deleted: " + user.Email)

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
