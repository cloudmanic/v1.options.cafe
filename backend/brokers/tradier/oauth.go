//
// Date: 9/6/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
)

type TradierAuth struct {
	DB models.Datastore
}

var (
	genericError = []byte("Something went wrong while authorizing your account. Please try again or contact help@options.cafe. Sorry for the trouble.")
)

// Json response object
type tokenResponse struct {
	Token        string `json:"access_token"`
	ExpiresSec   int64  `json:"expires_in"`
	IssueDateStr string `json:"issued_at"`
	RefreshToken string `json:"refresh_token"`
	Scope        string
	Status       string
}

//
// Obtain an Authorization Code - http://localhost:7080/tradier/authorize?user=1
//
func (t *TradierAuth) DoAuthCode(c *gin.Context) {

	// Make sure we have a user id.
	userId := c.Query("user")

	if userId == "" {
		services.Info(errors.New("Tradier - DoAuthCode - No user id provided."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	if userId == "" {
		services.Info(errors.New("Tradier - DoAuthCode - No user id provided."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure this is a valid user.
	u, _ := strconv.ParseUint(userId, 10, 32)
	user, err := t.DB.GetUserById(uint(u))

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure we have a broker_id.
	brokerId := c.Query("broker_id")

	if brokerId == "" {
		services.Info(errors.New("Tradier - DoAuthCode - No broker id provided."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	if brokerId == "" {
		services.Info(errors.New("Tradier - DoAuthCode - No broker id provided."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure this is a valid broker_id.
	broker := models.Broker{}

	if t.DB.New().Where("user_id = ? AND name = ? AND id = ?", user.Id, "Tradier", brokerId).Find(&broker).RecordNotFound() {
		services.Info(errors.New("Tradier - DoAuthCode - No broker found."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Log
	services.InfoMsg("Tradier authorization starting for " + user.Email)

	// Redirect to tradier to auth
	var url = apiBaseUrl + "/oauth/authorize?client_id=" + os.Getenv("TRADIER_CONSUMER_KEY") + "&scope=read,write,market,trade,stream&state=" + strconv.Itoa(int((broker.Id)))
	c.Redirect(302, url)
}

//
// Do Obtain an Authorization Code Callback - http://localhost:7080/tradier/callback
//
func (t *TradierAuth) DoAuthCallback(c *gin.Context) {
	// Make sure we have a code.
	code := c.Query("code")

	if code == "" {
		services.Info(errors.New("Tradier - DoAuthCallback - No auth code provided. (#1)"))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure we have a state.
	state := c.Query("state")

	if state == "" {
		services.Info(errors.New("Tradier - DoAuthCallback - No auth code provided. (#2)"))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Get broker from state
	broker := models.Broker{}

	if t.DB.New().Where("name = ? AND id = ?", "Tradier", state).Find(&broker).RecordNotFound() {
		services.Info(errors.New("Tradier - DoAuthCallback - No broker found."))
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	params := url.Values{}
	params.Set("scope", "read,write,trade,market,stream")
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", apiBaseUrl+"/oauth/accesstoken", body)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Headers
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(os.Getenv("TRADIER_CONSUMER_KEY"), os.Getenv("TRADIER_CONSUMER_SECRET"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	defer resp.Body.Close()

	// Make sure we got a good status code
	if resp.StatusCode != http.StatusOK {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Get the json out of the body.
	jsonBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Put json into an object.
	var tr tokenResponse

	err = json.Unmarshal(jsonBody, &tr)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure this request was approved.
	if tr.Status != "approved" {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Make sure this is a valid user.
	user, err := t.DB.GetUserById(broker.UserId)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": genericError})
		return
	}

	// Update Broker
	//broker.Status = "Active"
	broker.AccessToken = tr.Token
	broker.RefreshToken = tr.RefreshToken
	broker.TokenExpirationDate = time.Now().Add(time.Duration(tr.ExpiresSec) * time.Second).UTC()
	t.DB.UpdateBroker(broker)

	// We bootstrap the broker before setting it active this way the normal poller
	// does not overlap. NOTE: We set the broker to active in KickStartBroker that is
	// when the everyday poller kicks in.
	go t.DB.KickStartBroker(user, broker)

	// Just give the data building a bit of a head start
	time.Sleep(time.Second * 2)

	// Log
	services.InfoMsg("Tradier authorization completed for " + user.Email)

	// Return success redirect
	c.Redirect(302, os.Getenv("SITE_URL"))
}

//
// Check to see if we need to refresh the refresh token.
//
func (t *Api) DoRefreshAccessTokenIfNeeded(user models.User) error {
	// Get the different tradier brokers.
	brokers, err := t.DB.GetBrokerTypeAndUserId(user.Id, "Tradier")

	if err != nil {
		services.Info(err)
		return err
	}

	// Loop through and deal with each tradier broker in the db.
	for i := range brokers {

		// Skip over disabled brokers
		if brokers[i].Status == "Disabled" {
			continue
		}

		// Is it time to refresh
		if time.Now().UTC().Add(1 * time.Hour).After(brokers[i].TokenExpirationDate.UTC()) {

			err, msg := t.DoRefreshAccessToken(brokers[i])

			if err == nil {

				// Update the access token.
				t.muApiKey.Lock()
				t.ApiKey = msg
				t.muApiKey.Unlock()

				services.InfoMsg("Refreshed Tradier token : " + user.Email)

			} else {
				services.InfoMsg("Refreshed Tradier error : " + msg)
				services.Info(err)
			}

		}

	}

	// All done no errors
	return nil
}

//
// Get a new access token via the refresh token.
//
func (t *Api) DoRefreshAccessToken(broker models.Broker) (error, string) {
	// Decrypt the refresh token
	decryptRefreshToken, err := helpers.Decrypt(broker.RefreshToken)

	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Unable to decrypt message (#1)"
	}

	// Request and get an access token.
	data := strings.NewReader("grant_type=refresh_token&refresh_token=" + decryptRefreshToken)

	req, err := http.NewRequest("POST", apiBaseUrl+"/oauth/refreshtoken", data)

	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#1)"
	}

	req.SetBasicAuth(os.Getenv("TRADIER_CONSUMER_KEY"), os.Getenv("TRADIER_CONSUMER_SECRET"))
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#2)"
	}

	defer resp.Body.Close()

	// Get the json out of the body.
	jsonBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#4)"
	}

	// Make sure we got a good status code
	if resp.StatusCode != http.StatusOK {
		return errors.New("Tradier - DoRefreshAccessToken - Failed to get access token. (#3)"), string(jsonBody)
	}

	// Put json into an object.
	var tr tokenResponse

	err = json.Unmarshal(jsonBody, &tr)

	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#5)"
	}

	// Make sure this request was approved.
	if tr.Status != "approved" {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#6)"
	}

	// Update the database
	broker.AccessToken = tr.Token
	broker.RefreshToken = tr.RefreshToken
	broker.TokenExpirationDate = time.Now().Add(time.Duration(tr.ExpiresSec) * time.Second).UTC()
	t.DB.UpdateBroker(broker)

	// All done no errors
	return nil, tr.Token
}

/* End File */
