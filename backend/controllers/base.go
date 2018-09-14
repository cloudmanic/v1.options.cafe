//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"flag"
	"go/build"
	"net/http"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/users"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
	"github.com/gin-gonic/gin"
	env "github.com/jpfuentes2/go-env"
)

const defaultMysqlLimit = 100
const httpNoRecordFound = "No Record Found."
const httpGenericErrMsg = "Please contact support at help@options.cafe."

type Controller struct {
	DB                  models.Datastore
	UserActionChan      chan users.UserFeedAction
	WebsocketController *websocket.Controller
}

type ValidateRequest interface {
	Validate(models.Datastore) error
}

//
// Start up the controller.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")
}

//
// Add paging info to the response.
//
func (t *Controller) AddPagingInfoToHeaders(c *gin.Context, meta models.QueryMetaData) {
	c.Writer.Header().Set("X-Last-Page", strconv.FormatBool(meta.LastPage))
	c.Writer.Header().Set("X-Offset", strconv.Itoa(meta.Offset))
	c.Writer.Header().Set("X-Limit", strconv.Itoa(meta.Limit))
	c.Writer.Header().Set("X-No-Limit-Count", strconv.Itoa(meta.NoLimitCount))
}

//
// Validate and Create object.
//
func (t *Controller) ValidateRequest(c *gin.Context, obj ValidateRequest) error {

	// Bind the JSON that got sent into an object and validate.
	if err := c.ShouldBindJSON(obj); err == nil {

		// Run validation
		err := obj.Validate(t.DB)

		// If we had validation errors return them and do no more.
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return err
		}

	} else {
		services.Warning(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in body. There is a chance the JSON maybe valid but does not match the data type requirements. For example maybe you passed a string in for an integer."})
		return err
	}

	return nil
}

//
// Get / Set standard query parms
//
func GetSetPagingParms(c *gin.Context) (int, int, int) {
	// Convert page to int.
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	// We do not allow limits over defaultMysqlLimit
	if limit > defaultMysqlLimit {
		limit = defaultMysqlLimit
	}

	if limit == 0 {
		limit = defaultMysqlLimit
	}

	// Offset can't be less than 0
	if offset < 0 {
		offset = 0
	}

	// Page can't be less than 1
	if page < 1 {
		page = 1
	}

	// Return happy.
	return page, limit, offset
}

//
// Respond with an error or object. When we create a new object in the system
//
func (t *Controller) RespondCreated(c *gin.Context, payload interface{}, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, payload)
	}
}

//
// RespondJSON makes the response with payload as json format.
// This is used when we want the json back (used in websockets).
// If you do not need the json back just use c.JSON()
//
func (t *Controller) RespondJSON(c *gin.Context, status int, payload interface{}) string {

	response, err := json.Marshal(payload)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return ""
	}

	// Return json.
	c.JSON(200, payload)

	// We return the raw JSON
	return string(response)
}

//
// Return error.
//
func (t *Controller) RespondError(c *gin.Context, err error, msg string) bool {

	if err != nil {
		if flag.Lookup("test.v") == nil {
			services.Warning(err)
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return true
	}

	// No error.
	return false
}

/* End File */
