//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"app.options.cafe/brokers/tradier"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// Return groups in our database.
//
func (t *Controller) GetBrokers(c *gin.Context) {

	var results = []models.Broker{}

	// Run the query
	err := t.DB.Query(&results, models.QueryParam{
		UserId:   c.MustGet("userId").(uint),
		Limit:    defaultMysqlLimit,
		PreLoads: []string{"BrokerAccounts"},
	})

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	c.JSON(200, results)
}

//
// Create new broker. We default to disabled. The broker
// is enabled when the oauth handshake happens.
//
func (t *Controller) CreateBroker(c *gin.Context) {

	// Get User Id
	userId := c.MustGet("userId").(uint)

	// Setup Broker obj
	o := models.Broker{UserId: userId}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get inputs. We do it this way as we do not want to accept any other input for security reasons.
	name := gjson.Get(string(body), "name").String()
	displayName := gjson.Get(string(body), "display_name").String()

	// Validate name
	if len(name) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name field can not be empty."})
		return
	}

	// For now we only support tradier
	h := []string{"Tradier"}
	found, _ := helpers.InArray(name, h)

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Broker name not valid."})
		return
	}

	// Validate display name
	if len(displayName) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display Name field can not be empty."})
		return
	}

	// Load up data
	o.Name = name
	o.Status = "Disabled"
	o.DisplayName = displayName

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Update Sendy with this new fact.
	go services.SendyUnsubscribe("no-brokers", user.Email)
	go services.SendySubscribe("subscribers", user.Email, user.FirstName, user.LastName, "No", name, "", "No")

	// Create Screen
	err = t.DB.CreateNewRecord(&o, models.InsertParam{})
	t.RespondCreated(c, o, err)
}

//
// Update broker.
//
func (t *Controller) UpdateBroker(c *gin.Context) {

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get inputs. We do it this way as we do not want to accept any other input for security reasons.
	displayName := gjson.Get(string(body), "display_name").String()

	// Validate display name
	if len(displayName) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display Name field can not be empty."})
		return
	}

	// Setup Broker obj
	broker, err := t.DB.GetBrokerById(uint(id))

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Update data
	broker.DisplayName = displayName

	// Update broker
	t.DB.New().Save(&broker)

	// Return happy JSON
	c.JSON(200, broker)
}

//
// Make an API call to the broker and get balances.
//
func (t *Controller) GetBalances(c *gin.Context) {

	var apiKey string = ""
	var brokers = []models.Broker{}

	// Run the query to get brokers
	err := t.DB.Query(&brokers, models.QueryParam{
		UserId: c.MustGet("userId").(uint),
		Wheres: []models.KeyValue{
			{Key: "name", Value: "Tradier"},
			{Key: "Status", Value: "Active"},
		},
	})

	// Loop through the different brokers- TODO: This only supports one broker. We need to get balance from all brokers and merge data together.
	if !strings.HasSuffix(os.Args[0], ".test") {
		for _, row := range brokers {

			// Decrypt the access token
			_apiKey, err := helpers.Decrypt(row.AccessToken)

			if err != nil {
				t.RespondError(c, err, httpGenericErrMsg)
				return
			}

			apiKey = _apiKey
		}
	}

	// // Figure out which broker connection to setup.
	// switch broker.Name {

	// case "Tradier":
	// 	brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: false}

	// case "Tradier Sandbox":
	// 	brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: true}

	// default:
	// 	services.InfoMsg("Order: Unknown Broker : " + broker.Name)

	// }

	// Setup the broker
	broker := tradier.Api{
		DB:     t.DB,
		ApiKey: apiKey,
	}

	// Make API call to broker.
	result, err := broker.GetBalances()

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Return happy JSON
	c.JSON(200, result)
}

/* End File */
