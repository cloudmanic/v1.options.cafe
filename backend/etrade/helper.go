//TODO: Check for API errors returned at API level
package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/mrjones/oauth"
)

var (
	baseURL    = "https://api.etrade.com/v1/"
	sandboxURL = "https://api.etrade.com/%s/rest/%s"
	jsonURL    = ".json"
)

type ETradeClient struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
	url         string
}

type strURLParameter []string

// Creates a comma separated string
func (q strURLParameter) String() string {
	return strings.Join(q, ",")
}

type IntDollar int64

func (i IntDollar) String() (s string) {
	dollars := i / 100
	cents := i % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}

func New(consumerID, consumerSecret string) (client ETradeClient, err error) {

	c := oauth.NewConsumer(
		consumerID,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.etrade.com/oauth/request_token",
			AuthorizeTokenUrl: "https://us.etrade.com/e/t/etws/authorize",
			AccessTokenUrl:    "https://api.etrade.com/oauth/access_token",
		})
	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		return client, err
	}
	url = fmt.Sprintf("https://us.etrade.com/e/t/etws/authorize?key=%s&token=%s", consumerID, requestToken.Token)
	fmt.Println("(1) Go to: " + url + "")
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		return client, err
	}

	client.consumer = c
	client.accessToken = accessToken
	client.url = baseURL
	return
}

func (client *ETradeClient) SetToSandBox() {
	client.url = sandboxURL
}

func (client *ETradeClient) SetToProduction() {
	client.url = baseURL
}

//
// makeRequest sends and API request with oauth1 signature.
// We return a raw JSON string.
//
func (client *ETradeClient) makeRequest(requestURL string, params map[string]string) (raw string, err error) {
	// Build oauth consumer
	r, err := client.consumer.Get(requestURL, params, client.accessToken)

	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	// Check http status.
	if r.StatusCode == http.StatusNotModified {
		return raw, nil
	} else if r.StatusCode != http.StatusOK {
		return raw, fmt.Errorf("Status code not valid for request: %s", r.Status)
	}

	// Read body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	// Return raw JSON
	return string(body), err
}

func convertToIntDollar(v interface{}) (i IntDollar, err error) {
	switch v := v.(type) {
	default:
		return i, fmt.Errorf("Unexpected type in response for Value %v", v)
	case float32:
		i = IntDollar(math.Floor((float64(v) * 100) + .5))
	case string:
		//This doesn't handle strings without a decimal correctly
		//Need to unit test more
		v = strings.Replace(v, ".", "", -1)
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return i, err
		}
		i = IntDollar(value)
	case float64:
		i = IntDollar(math.Floor((v * 100) + .5))
	}
	return
}
