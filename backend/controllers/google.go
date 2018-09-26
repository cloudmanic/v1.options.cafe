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

	state := randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	// Redirect to google to login.
	c.Redirect(302, googleConf.AuthCodeURL(state))
}

//
// Deal with the call back from Google with the code.
//
func (t *Controller) DoGoogleCallback(c *gin.Context) {

	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")

	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
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

  // Redirect back to main site
  c.Redirect(302, os.Getenv("SITE_URL"))
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
