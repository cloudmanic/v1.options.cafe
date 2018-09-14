package tradier

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
)

const (
	apiBaseUrl  = "https://api.tradier.com/v1"
	sandBaseUrl = "https://sandbox.tradier.com/v1"
)

type Api struct {
	DB      models.Datastore
	Sandbox bool

	muActiveSymbols sync.Mutex
	activeSymbols   string

	muApiKey sync.Mutex
	ApiKey   string
}

//
// Start up the controller.
//
func init() {
	// Helpful for testing
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")
}

//
// Add symbols to active list.
//
func (t *Api) SetActiveSymbols(symbols []string) {

	// Lock da memory
	t.muActiveSymbols.Lock()
	defer t.muActiveSymbols.Unlock()

	t.activeSymbols = strings.Join(symbols, ",")

}

//
// Send a POST request to Tradier. Returns the JSON string or an error
//
func (t *Api) SendPostRequest(urlStr string, params url.Values) (string, error) {

	// Setup http client
	client := &http.Client{}

	// Get url to api
	apiUrl := apiBaseUrl

	if t.Sandbox {
		apiUrl = sandBaseUrl
	}

	// Setup api request
	req, _ := http.NewRequest("POST", apiUrl+urlStr, bytes.NewBufferString(params.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	// Close Body
	defer res.Body.Close()

	// Read the data we got.
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprint("API did not return 200, It returned (", apiUrl+urlStr, ")", res.StatusCode, " ", string(body)))
	}

	// Return happy.
	return string(body), nil
}

//
// Send a GET request to Tradier. Returns the JSON string or an error
//
func (t *Api) SendGetRequest(urlStr string) (string, error) {

	// Setup http client
	client := &http.Client{}

	// Get url to api
	apiUrl := apiBaseUrl

	if t.Sandbox {
		apiUrl = sandBaseUrl
	}

	// Setup api request
	req, _ := http.NewRequest("GET", apiUrl+urlStr, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	// Close Body
	defer res.Body.Close()

	// Read the data we got.
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprint("API did not return 200, It returned (", apiUrl+urlStr, ")", res.StatusCode, " ", string(body)))
	}

	// Return happy.
	return string(body), nil
}

/* End File */
