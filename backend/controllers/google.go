//
// Date: 9/25/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConf *oauth2.Config

//
// Init.
//
func init() {

	// Setup google Oauth
	googleConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BACKEND_URL") + "/oauth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}

}

//
// Do Google Auth login
//
func (t *Controller) DoGoogleAuthLogin(c *gin.Context) {

	// Make sure we pass in a shared key to be used later to allow access to an access token
	sharedKey := c.Query("shared")

	if len(sharedKey) <= 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid shared key: %s", c.Query("shared")))
		return
	}

	// Set state
	state := randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("google_auth_shared_key", sharedKey)
	session.Save()

	// Redirect to google to login.
	c.Redirect(302, googleConf.AuthCodeURL(state))
}

//
// Deal with the call back from Google with the code.
// If success we create a shared secret between the front-end
// and the backkend. When we started this process the front end pased in
// googleAuthSharedKey as a shared key. This key is stored in the browser's
// local storage. Upon success auth with google we create a new random key.
// we put these two keys together as a MD5 hash and store in redis. We then redirect
// back to the browser with the server's shared key. The server then sends back both shared
// keys to obtain an access token.
//
func (t *Controller) DoGoogleCallback(c *gin.Context) {

	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	googleAuthSharedKey := session.Get("google_auth_shared_key").(string)

	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	if len(googleAuthSharedKey) <= 0 {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("No google auth shared key found"))
		return
	}

	// Exchange the code for an access token.
	tok, err := googleConf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Now that we have an access token get the user profile.
	client := googleConf.Client(oauth2.NoContext, tok)

	profile, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer profile.Body.Close()

	data, err := ioutil.ReadAll(profile.Body)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get data google returned.
	subId := gjson.Get(string(data), "sub").String()
	email := gjson.Get(string(data), "email").String()
	firstName := gjson.Get(string(data), "given_name").String()
	lastName := gjson.Get(string(data), "family_name").String()

	fmt.Println(subId, email, firstName, lastName)

	// Shared.
	googleShared := randToken()

	// Encrypt the string we are storing.
	hash, err := helpers.Encrypt(googleAuthSharedKey + ":" + googleShared)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Create md5 hash of the string
	md5Hash := helpers.GetMd5(googleAuthSharedKey + ":" + googleShared)

	// Store another shared secret to help give google auth an access token.
	cache.SetExpire("google_auth_shared_key_"+md5Hash, (time.Minute * 1), hash)

	// Redirect back to main site
	c.Redirect(302, os.Getenv("SITE_URL")+"?google_auth_success="+googleShared)
}

//
// Get random token
//
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

/* End File */
