package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"io/ioutil"
	"net/http"
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
	body, _ := ioutil.ReadAll(res.Body)

	// Bust open the Profile.
	var ws map[string]types.TmpUserProfile

	if err := json.Unmarshal(body, &ws); err != nil {
		return user, err
	}

	// Set the object.
	user.Id = ws["profile"].Id
	user.Name = ws["profile"].Name

	// Loop through and set accounts.
	for _, row := range ws["profile"].Accounts {

		user.Accounts = append(user.Accounts, types.UserProfileAccounts{
			AccountNumber:  row.AccountNumber,
			Classification: row.Classification,
			DayTrader:      row.DayTrader,
			OptionLevel:    row.OptionLevel,
			Status:         row.Status,
			Type:           row.Type})

	}

	// Return happy
	return user, nil

}

/* End File */
