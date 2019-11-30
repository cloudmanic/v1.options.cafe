package main

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// Account struct
type Account struct {
	ID              string `json:"accountId"`
	IDKey           string `json:"accountIdKey"`
	Name            string `json:"accountName"`
	Description     string `json:"accountDesc"`
	Mode            string `json:"accountMode"`
	Type            string `json:"accountType"`
	InstitutionType string `json:"institutionType"`
	Status          string `json:"accountStatus"`
	ClosedDate      int64  `json:"closedDate"`
	Balance         Balance
}

// Balance struct
type Balance struct {
	ID                string
	TotalAccountValue float64
}

//
// ListAccounts will return a list of all accounts related to this etrade user.
//
func (client ETradeClient) ListAccounts() (accounts []Account, err error) {
	// Build url
	url := client.url + "accounts/list.json"

	// Make request.
	jsonRt, err := client.makeRequest(url, nil)

	if err != nil {
		return accounts, err
	}

	// Get to the guts.
	result := gjson.Get(jsonRt, "AccountListResponse.Accounts.Account")

	// Loop  through and build out the structs
	for _, row := range result.Array() {
		a := Account{}

		if err := json.Unmarshal([]byte(row.String()), &a); err != nil {
			return []Account{}, err
		}

		// We do not care about closed accounts.
		if a.Status == "CLOSED" {
			continue
		}

		// Get the balance for this account.
		balance, err := client.ListAccountBalance(a.IDKey)

		if err != nil {
			return []Account{}, err
		}

		// Add in balance
		a.Balance = balance

		// Add to the array
		accounts = append(accounts, a)
	}

	// Return happy
	return accounts, nil
}

//
// ListAccountBalance will return the balance information for a particular account by IDKey
//
func (client ETradeClient) ListAccountBalance(idKey string) (balance Balance, err error) {
	// Build url
	url := fmt.Sprintf("%saccounts/%s/balance.json", client.url, idKey)

	// Get GET Params
	p := map[string]string{}
	p["realTimeNAV"] = "true"
	p["instType"] = "BROKERAGE"

	// Make request.
	jsonRt, err := client.makeRequest(url, p)

	if err != nil {
		return balance, err
	}

	// Build object
	balance.ID = gjson.Get(jsonRt, "BalanceResponse.accountId").String()
	balance.TotalAccountValue = gjson.Get(jsonRt, "BalanceResponse.Computed.RealTimeValues.totalAccountValue").Float()

	// Return happy
	return balance, nil
}

/* End File */
