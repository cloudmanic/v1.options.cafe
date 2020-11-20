package tradier

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"app.options.cafe/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get user profile
//
func (t *Api) GetUserProfile() (types.UserProfile, error) {

	var user types.UserProfile

	// Setup http client
	client := &http.Client{}

	// Get url to api
	apiUrl := apiBaseUrl

	if t.Sandbox {
		apiUrl = sandBaseUrl
	}

	// Setup api request
	req, _ := http.NewRequest("GET", apiUrl+"/user/profile", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return user, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return user, errors.New(fmt.Sprint("GetUserProfile API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	jsonResponse, _ := ioutil.ReadAll(res.Body)

	// Set the object.
	user.Id = gjson.Get(string(jsonResponse), "profile.id").String()
	user.Name = gjson.Get(string(jsonResponse), "profile.name").String()

	// See if we have one account or many account
	vo := gjson.Get(string(jsonResponse), "profile.account.account_number")

	// Only one account
	if vo.Exists() {

		// Add one account.
		user.Accounts = append(user.Accounts, types.UserProfileAccounts{
			AccountNumber:  gjson.Get(string(jsonResponse), "profile.account.account_number").String(),
			Classification: gjson.Get(string(jsonResponse), "profile.account.classification").String(),
			DayTrader:      gjson.Get(string(jsonResponse), "profile.account.day_trader").Bool(),
			OptionLevel:    int(gjson.Get(string(jsonResponse), "profile.account.option_level").Int()),
			Status:         gjson.Get(string(jsonResponse), "profile.account.status").String(),
			Type:           gjson.Get(string(jsonResponse), "profile.account.type").String(),
		})

	} else {

		// Loop through the accounts because there is more than one.
		vo1 := gjson.Get(string(jsonResponse), "profile.account")

		for _, row := range vo1.Array() {

			user.Accounts = append(user.Accounts, types.UserProfileAccounts{
				AccountNumber:  gjson.Get(row.String(), "account_number").String(),
				Classification: gjson.Get(row.String(), "classification").String(),
				DayTrader:      gjson.Get(row.String(), "day_trader").Bool(),
				OptionLevel:    int(gjson.Get(row.String(), "option_level").Int()),
				Status:         gjson.Get(row.String(), "status").String(),
				Type:           gjson.Get(row.String(), "type").String(),
			})

		}

	}

	// Return happy
	return user, nil

}

/* End File */
