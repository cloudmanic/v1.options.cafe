package tradier

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/app.options.cafe/backend/models"
)

const (
	apiBaseUrl = "https://api.tradier.com/v1"
)

type Api struct {
	DB              models.Datastore
	muActiveSymbols sync.Mutex
	activeSymbols   string

	muApiKey sync.Mutex
	ApiKey   string
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
// Send a GET request to Tradier. Returns the JSON string or an error
//
func (t *Api) SendGetRequest(urlStr string) (string, error) {

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", apiBaseUrl+urlStr, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprint("API did not return 200, It returned (", urlStr, ")", res.StatusCode))
	}

	// Read the data we got.
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	// Return happy.
	return string(body), nil
}

/* End File */
