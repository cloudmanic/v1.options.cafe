//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
)

//
// Get account invoice history
//
func (t *Controller) BillingHistory(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Get the history
	history, err := t.DB.GetInvoiceHistoryWithStripe(user)

	if t.RespondError(c, err, "Invoice not found. Please contact help@options.cafe") {
		return
	}

	// Return happy
	c.JSON(200, history)
}

//
// Subscribe to a plan.
//
func (t *Controller) SubscribeUser(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Read data from POST request.
	plan := gjson.Get(string(body), "plan").String()
	token := gjson.Get(string(body), "token").String()
	coupon := gjson.Get(string(body), "coupon").String()

	// If plan is monthly
	if plan == "monthly" {
		plan = os.Getenv("STRIPE_MONTHLY_PLAN")
	}

	// If plan is yearly
	if plan == "yearly" {
		plan = os.Getenv("STRIPE_YEARLY_PLAN")
	}

	// Talk to stripe and setup the account.
	err = t.DB.CreateNewUserWithStripe(user, plan, token, coupon)

	if t.RespondError(c, err, "Unable to upgrade your account. Please contact help@options.cafe") {
		return
	}

	// Return happy
	c.JSON(202, nil)

}

//
// Verify coupon
//
func (t *Controller) VerifyCoupon(c *gin.Context) {

	// Get code from Stripe
	coupon, err := services.StripeGetCoupon(c.Param("code"))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Setup return value
	type RtStruct struct {
		Valid      bool    `json:"valid"`
		Name       string  `json:"name"`
		Code       string  `json:"code"`
		AmountOff  int64   `json:"amount_off"`
		PercentOff float64 `json:"percent_off"`
		Duration   string  `json:"duration"`
	}

	rt := RtStruct{
		Valid:      coupon.Valid,
		Name:       coupon.Name,
		Code:       coupon.ID,
		AmountOff:  coupon.AmountOff,
		PercentOff: coupon.PercentOff,
		Duration:   string(coupon.Duration),
	}

	// Return happy JSON
	c.JSON(200, rt)
}

//
// Apply a coupon (discount) to the account.
//
func (t *Controller) ApplyCoupon(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Read data from POST request.
	couponCode := gjson.Get(string(body), "coupon_code").String()

	// Add the credit card to stripe
	err = t.DB.ApplyCoupon(user, couponCode)

	if t.RespondError(c, err, "Unable to add coupon to your account. Please contact help@options.cafe") {
		return
	}

	// Return happy
	c.JSON(202, nil)
}

//
// Add credit card to the account. If one is already on the account we
// replace the card and add the new one. We pass in a stripe token.
//
func (t *Controller) UpdateCreditCard(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Read data from PUT request.
	token := gjson.Get(string(body), "token").String()

	// Add the credit card to stripe
	err = t.DB.UpdateCreditCard(user, token)

	if t.RespondError(c, err, "Unable to add credit card to your account. Please contact help@options.cafe") {
		return
	}

	// Do we have a coupon code?
	couponCode := gjson.Get(string(body), "coupon_code").String()

	if len(couponCode) > 0 {
		services.Info("Coupon code passed with credit card token: " + couponCode + " - " + user.Email)
		err = t.DB.ApplyCoupon(user, couponCode)

		if err != nil {
			services.BetterError(err)
		}
	}

	// Return happy
	c.JSON(202, nil)
}

//
// Get Me. Return user profile.
//
func (t *Controller) GetProfile(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Return happy JSON
	c.JSON(200, user)
}

//
// Update Me. Update a user profile.
//
func (t *Controller) UpdateProfile(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Setup BrokerAccount obj
	o := models.User{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// We only allow a few fields to be updated via the API
	user.Email = o.Email
	user.FirstName = o.FirstName
	user.LastName = o.LastName
	user.Phone = o.Phone
	user.Address = o.Address
	user.City = o.City
	user.State = o.State
	user.Country = o.Country
	user.Zip = o.Zip

	// Update BrokerAccount
	t.DB.New().Save(&user)

	// Return happy JSON
	c.JSON(202, user)
}

//
// Get the current subscription details.
//
func (t *Controller) GetSubscription(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Build a default subscription
	sub := models.UserSubscription{
		TrialDays:  helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT")),
		Status:     "trialing",
		Started:    user.CreatedAt,
		TrialStart: user.CreatedAt,
		TrialEnd:   user.TrialExpire,
	}

	// See if we have a subscription
	if len(user.StripeSubscription) > 0 {

		// Get subscription with stripe
		sub2, err := t.DB.GetSubscriptionWithStripe(user)

		if t.RespondError(c, err, "Subscription not found. Please contact help@options.cafe") {
			return
		}

		sub = sub2

	}

	// Return happy JSON
	c.JSON(200, sub)
}

//
// Reset the users password
//
func (t *Controller) ResetPassword(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	newPass := gjson.Get(string(body), "new_password").String()
	currentPass := gjson.Get(string(body), "current_password").String()

	// Now that we know the user lets make sure the password that was posted in was at least 6 chars.
	err = t.DB.ValidatePassword(newPass)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a password at least 6 chars long."})
		return
	}

	// Validate password here by comparing hashes nil means success
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPass))

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect current password."})
		return
	}

	// Change password
	err = t.DB.ResetUserPassword(user.Id, newPass)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to reset password."})
		return
	}

	// Return happy
	c.JSON(202, nil)
}

/* End File */
