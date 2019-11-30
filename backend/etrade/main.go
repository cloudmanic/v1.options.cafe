package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/mrjones/oauth"
)

var consumerID = ""
var consumerSecret = ""

var accessToken = ""
var accessSecret = ""

func main() {

	// client, err := New(consumerID, consumerSecret)
	//
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// spew.Dump(client.accessToken)

	c := oauth.NewConsumer(
		consumerID,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.etrade.com/oauth/request_token",
			AuthorizeTokenUrl: "https://us.etrade.com/e/t/etws/authorize",
			AccessTokenUrl:    "https://api.etrade.com/oauth/access_token",
		})

	at := oauth.AccessToken{
		Token:  accessToken,
		Secret: accessSecret,
	}

	client := ETradeClient{}
	client.consumer = c
	client.accessToken = &at
	client.url = baseURL

	accts, err := client.ListAccounts()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	spew.Dump(accts)

	// balance, err := client.ListAccountBalance("XXXXXXX")
	//
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// spew.Dump(balance)

	return
}
