package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"app.options.cafe/backend/library/realip"
	"app.options.cafe/backend/library/services"
)

// TODO: Lots of duplicate code in here with setting headers and such. Should clean up. Also see Login, and Register.

//
// Rest the password after they clicked on the email - Step #2
//
func (t *Controller) DoResetPassword(w http.ResponseWriter, r *http.Request) {

	// Manage OPTIONS requests
	if (os.Getenv("APP_ENV") == "local") && (r.Method == http.MethodOptions) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
		return
	}

	// // Make sure this is a post request.
	// if r.Method == http.MethodGet {
	// 	t.HtmlMainTemplate(w, r)
	// 	return
	// }

	// Make sure this is a post request.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Content-Type", "application/json")

	// Decode json passed in
	decoder := json.NewDecoder(r.Body)

	type ResetPost struct {
		Hash     string
		Password string
	}

	var post ResetPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Error(err, "DoResetPassword - Failed to decode JSON posted in")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Something went wrong while setting up your forgot password request. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))
		return
	}

	defer r.Body.Close()

	// Get the user based on the hash we passed in.
	user, err := t.DB.GetUserFromToken(post.Hash)

	if err != nil {
		services.Error(err, "GetUserFromToken - Unable to reset password.")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Sorry, It seems your reset token has expired.\"}"))

		return
	}

	// Now that we know the user lets make sure the password that was posted in was at least 6 chars.
	err = t.DB.ValidatePassword(post.Password)

	if err != nil {
		services.Error(err, "ValidatePassword - User posted in a password that was less than 6 chars (this should not happen, angular validates).")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Please enter a password at least 6 chars long.\"}"))

		return
	}

	// Now that we know who the user is lets reset the users password.
	err = t.DB.ResetUserPassword(user.Id, post.Password)

	if err != nil {
		services.Error(err, "ResetUserPassword - Password reset failed.")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Something went wrong while changing your password. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))

		return
	}

	// Lastly delete the reset password hash.
	err = t.DB.DeleteForgotPasswordByToken(post.Hash)

	if err != nil {
		services.Error(err, "DeleteForgotPasswordByToken - Deleting reset password token")
	}

	// Build response
	type Response struct {
		Message string `json:"message"`
	}

	resObj := &Response{Message: "Success! Your password has been updated."}

	resJson, err := json.Marshal(resObj)

	if err != nil {
		services.Error(err, "DoResetPassword - Unable to reset password. (json.Marshal)")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Something went wrong while changing your password. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))

		return
	}

	// Return success json.
	w.Write(resJson)

}

//
// Post back to setup a forgot password request. - Step #1 (send email request to reset)
//
func (t *Controller) DoForgotPassword(w http.ResponseWriter, r *http.Request) {

	// Manage OPTIONS requests
	if (os.Getenv("APP_ENV") == "local") && (r.Method == http.MethodOptions) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
		return
	}

	// // Make sure this is a post request.
	// if r.Method == http.MethodGet {
	// 	t.HtmlMainTemplate(w, r)
	// 	return
	// }

	// Make sure this is a post request.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Content-Type", "application/json")

	// Decode json passed in
	decoder := json.NewDecoder(r.Body)

	type ForgotPost struct {
		Email string
	}

	var post ForgotPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Error(err, "DoForgotPassword - Failed to decode JSON posted in")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Something went wrong while setting up your forgot password request. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))
		return
	}

	defer r.Body.Close()

	// Request a reset password request.
	err = t.DB.DoResetPassword(post.Email, realip.RealIP(r))

	if err != nil {
		services.Error(err, "DoForgotPassword - Unable to reset password.")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Sorry, we could not find your account.\"}"))

		return
	}

	type Response struct {
		Message string `json:"message"`
	}

	resObj := &Response{Message: "Success! Please check your email for next steps."}

	resJson, err := json.Marshal(resObj)

	if err != nil {
		services.Error(err, "DoForgotPassword - Unable to reset password. (json.Marshal)")

		// Respond with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Something went wrong while setting up your forgot password request. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))

		return
	}

	// Return success json.
	w.Write(resJson)

}
