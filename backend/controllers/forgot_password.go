package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/realip"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

//
// Rest the password after they clicked on the email - Step #2
//
func (t *Controller) DoResetPassword(c *gin.Context) {

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type ResetPost struct {
		Hash     string
		Password string
	}

	var post ResetPost

	err := decoder.Decode(&post)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Get the user based on the hash we passed in.
	user, err := t.DB.GetUserFromToken(post.Hash)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, It seems your reset token has expired."})
		return
	}

	// Now that we know the user lets make sure the password that was posted in was at least 6 chars.
	err = t.DB.ValidatePassword(post.Password)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a password at least 6 chars long."})
		return
	}

	// Now that we know who the user is lets reset the users password.
	err = t.DB.ResetUserPassword(user.Id, post.Password)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	// Lastly delete the reset password hash.
	err = t.DB.DeleteForgotPasswordByToken(post.Hash)

	if err != nil {
		services.BetterError(err)
	}

	// Build response
	type Response struct {
		Message string `json:"message"`
	}

	resObj := &Response{Message: "Success! Your password has been updated."}

	// Return success json.
	c.JSON(200, resObj)
}

//
// Post back to setup a forgot password request. - Step #1 (send email request to reset)
//
func (t *Controller) DoForgotPassword(c *gin.Context) {

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type ForgotPost struct {
		Email string
	}

	var post ForgotPost

	err := decoder.Decode(&post)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Request a reset password request.
	err = t.DB.DoResetPassword(post.Email, realip.RealIP(c.Request))

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not find your account."})
		return
	}

	type Response struct {
		Message string `json:"message"`
	}

	resObj := &Response{Message: "Success! Please check your email for next steps."}

	// Return success json.
	c.JSON(200, resObj)
}
