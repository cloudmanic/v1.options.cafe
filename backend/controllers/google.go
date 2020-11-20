//
// Date: 9/25/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"app.options.cafe/library/cache"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/realip"
	"app.options.cafe/library/services"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConf *oauth2.Config

type GoogleSessionStore struct {
	Type          string `json:"type"`
	UserId        uint   `json:"user_id"`
	SessionSecret string `json:"session_secret"`
	SubId         string `json:"sub_id"`
	Email         string `json:"email"`
	First         string `json:"first"`
	Last          string `json:"last"`
	Redirect      string `json:"redirect"`
}

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
// Start a Google login session. We return a shared key.
// This shared key is good for 1 min and will help authenticate
// a front-end app.
//
func (t *Controller) DoStartGoogleLoginSession(c *gin.Context) {

	// Posted data.
	body, _ := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	// Set shared key
	sessionRedirect := gjson.Get(string(body), "redirect").String()
	sessionType := gjson.Get(string(body), "type").String()
	sessionKey := gjson.Get(string(body), "session_key").String()
	sessionSecret := randToken()

	// Build json to store
	jsonStore, _ := json.Marshal(GoogleSessionStore{
		UserId:        0,
		Type:          sessionType,
		SessionSecret: sessionSecret,
		Redirect:      sessionRedirect,
	})

	// Encrypt the string we are storing.
	hash, _ := helpers.Encrypt(string(jsonStore))

	// Store session key in redis.
	cache.SetExpire("google_auth_session_"+sessionKey, (time.Minute * 1), hash)

	// Return success json.
	c.JSON(200, gin.H{"session_secret": sessionSecret})
}

//
// Do Google Auth login
//
func (t *Controller) DoGoogleAuthLogin(c *gin.Context) {

	// Make sure we pass in a shared key to be used later to allow access to an access token
	sessionKey := c.Query("session_key")

	if len(sessionKey) <= 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid session key (#003): %s", c.Query("session_key")))
		return
	}

	// Make sure this is a valid session key
	var temp interface{}
	found, err := cache.Get("google_auth_session_"+sessionKey, &temp)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !found {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid session key (#004): %s", c.Query("session_key")))
		return
	}

	// Redirect to google to login.
	c.Redirect(302, googleConf.AuthCodeURL(sessionKey))
}

//
// Deal with the call back from Google with the code.
//
func (t *Controller) DoGoogleCallback(c *gin.Context) {

	// Handle the exchange code to initiate a transport.
	sessionKey := c.Query("state")

	// Make sure this is a valid session key
	var temp interface{}
	found, err := cache.Get("google_auth_session_"+sessionKey, &temp)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !found {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid session key (#001): %s", sessionKey))
		return
	}

	// Decrypt json blob
	jsonRt, err := helpers.Decrypt(temp.(string))

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// JSON to vars
	sessionRedirect := gjson.Get(jsonRt, "redirect").String()
	sessionType := gjson.Get(jsonRt, "type").String()
	sessionSecret := gjson.Get(jsonRt, "session_secret").String()

	// Now connect to google to get a token. Exchange the code for an access token.
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

	// Log for later support issues.
	services.InfoMsg("Google auth user information: " + string(data))

	// Get data google returned.
	subId := gjson.Get(string(data), "sub").String()
	email := gjson.Get(string(data), "email").String()
	firstName := gjson.Get(string(data), "given_name").String()
	lastName := gjson.Get(string(data), "family_name").String()

	// Get user from db
	user, err := t.DB.GetUserByGoogleSubId(subId)

	if sessionType == "login" {

		if err != nil {
			// Redirect back as user was not found.
			c.Redirect(302, sessionRedirect+"?google_auth_failed=user-not-found")
			return
		}

	} else {

		// They should login instead of register.
		if err == nil {
			// Redirect back as user was not found.
			c.Redirect(302, sessionRedirect+"?google_auth_failed=user-already-in-system")
			return
		}

	}

	// Build json to store by adding in user id
	jsonStore, _ := json.Marshal(GoogleSessionStore{
		UserId:        user.Id,
		SessionSecret: sessionSecret,
		Type:          sessionType,
		SubId:         subId,
		Email:         email,
		First:         firstName,
		Last:          lastName,
		Redirect:      sessionRedirect,
	})

	// Encrypt the string we are storing.
	hash, err := helpers.Encrypt(string(jsonStore))

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Store session key in redis.
	cache.SetExpire("google_auth_session_"+sessionKey, (time.Minute * 1), hash)

	// Redirect back to main site
	c.Redirect(302, sessionRedirect+"?google_auth_success=true")
}

//
// Get an access token after a google auth.
//
func (t *Controller) DoGetAccessTokenAfterGoogleAuth(c *gin.Context) {

	// Posted data.
	body, _ := ioutil.ReadAll(c.Request.Body)

	// Set shared key
	sessionKey := gjson.Get(string(body), "session_key").String()
	sessionSecret := gjson.Get(string(body), "session_secret").String()
	clientId := gjson.Get(string(body), "client_id").String()
	grantType := gjson.Get(string(body), "grant_type").String()

	defer c.Request.Body.Close()

	// First we validate the grant type and client id. Make sure this is a known application.
	app, err := t.DB.ValidateClientIdGrantType(clientId, grantType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client_id or grant type."})
		return
	}

	// Get the google auth session from Redis
	var temp interface{}
	found, err := cache.Get("google_auth_session_"+sessionKey, &temp)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !found {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid session key (#002): %s", c.Query("session_key")))
		return
	}

	// Decrypt json blob
	jsonRt, err := helpers.Decrypt(temp.(string))

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Set shared key
	redisSession := GoogleSessionStore{}
	if err := json.Unmarshal([]byte(jsonRt), &redisSession); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Verify the session secret
	if redisSession.SessionSecret != sessionSecret {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid session secret."))
		return
	}

	// Clear cache
	cache.SetExpire("google_auth_session_"+sessionKey, (time.Second * 1), "")

	// --- If we made it this far we can issue an access token.

	// Login user in by id
	user, err := t.DB.LoginUserById(uint(redisSession.UserId), app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Info(err)

		// Are we logging in or registering
		if redisSession.Type == "login" {

			// Respond with error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not find your account."})
			return

		} else {

			// Register new account
			userNew, err := t.DB.CreateUserFromGoogle(redisSession.First, redisSession.Last, redisSession.Email, redisSession.SubId, app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

			if err != nil {

				// Respond with error
				services.Info(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not register your account."})
				return

			}

			// Set the user.
			user = userNew
		}
	}

	// Return success json.
	c.JSON(200, gin.H{"access_token": user.Session.AccessToken, "user_id": user.Id, "broker_count": len(user.Brokers), "token_type": "bearer"})
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
