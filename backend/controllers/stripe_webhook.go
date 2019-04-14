//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
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
	services.InfoMsg("Stripe Webhook Received : " + event.Type + " - " + event.ID)

	// If there is no customer value there is nothing we need to do.
	// As of now all the events we care about have a customer attached to them.
	if len(event.GetObjectValue("customer")) == 0 {
		services.InfoMsg("Stripe no customer data found (this is expected).")
		c.JSON(200, gin.H{"status": "success"})
		return
	}

	// Figure out what user this event is for.
	user, err := t.DB.GetUserByStripeCustomer(event.GetObjectValue("customer"))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Figure out how we handle this event.
	switch event.Type {

	// customer.subscription.deleted
	case "customer.subscription.deleted":
		t.StripeEventSubscriptionDeleted(user)

	// invoice.payment_failed
	case "invoice.payment_failed":
		t.StripeEventPaymentFailed(user)

	// charge.succeeded
	case "charge.succeeded":
		t.StripeEventChargeSucceeded(user, event.GetObjectValue("balance_transaction"))

	}

	// Tell stripe all went well.
	c.JSON(200, gin.H{"status": "success"})
}

// -------------------- Stripe Events --------------------------- //

//
// Stripe event : charge.succeeded
//
func (t *Controller) StripeEventChargeSucceeded(user models.User, balanceTransaction string) {

	// Just make sure the user is in Active state
	user.Status = "Active"

	// Update the user's state
	t.DB.UpdateUser(&user)

	// Log event.
	services.InfoMsg("Stripe Subscription Update State Changed To " + user.Status + " : " + user.Email)

	// Get the amount we received and fee. This is something we could later send into Skyclerk
	tb, err := services.StripeGetBalanceTransaction(balanceTransaction)

	if err != nil {
		services.Info(err)
		return
	}

	// Log BalanceTransaction.
	fee := float32(float32(tb.Fee) / 100)
	amount := float32(float32(tb.Amount) / 100)

	services.InfoMsg("Stripe Payment Received From " + user.Email + ": Amount $" + fmt.Sprintf("%f", amount) + " Fee: $" + fmt.Sprintf("%f", fee))
}

//
// Stripe event : invoice.payment_failed
//
func (t *Controller) StripeEventPaymentFailed(user models.User) {

	// Delinquent means they are past due.
	user.Status = "Delinquent"

	// Update the user's state
	t.DB.UpdateUser(&user)

	// Log event.
	services.InfoMsg("Stripe Subscription Update State Changed To " + user.Status + " : " + user.Email)

	// Tell slack about this.
	go services.SlackNotify("#events", "Options Cafe User State Changed To "+user.Status+" : "+user.Email)
}

//
// Stripe event : customer.subscription.deleted
//
func (t *Controller) StripeEventSubscriptionDeleted(user models.User) {

	// Remove the subscription from the database.
	user.Status = "Delinquent"
	user.StripeSubscription = ""
	t.DB.UpdateUser(&user)

	// Log event.
	services.InfoMsg("Stripe Subscription Deleted: " + user.Email)

	// // Send email telling the user this happened.
	// var url = "https://app.options.cafe"

	// go email.Send(
	// 	user.Email,
	// 	"Options Cafe : Subscription Canceled",
	// 	emails.GetSubscriptionCanceledHtml(user.FirstName, url),
	// 	emails.GetSubscriptionCanceledText(user.FirstName, url))

	// Tell slack about this.
	go services.SlackNotify("#events", "Options Cafe Subscription StripeEventSubscriptionDeleted : "+user.Email)

}

/* End File */
