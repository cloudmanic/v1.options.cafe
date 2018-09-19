//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"encoding/json"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get Positions
//
func (t *Api) GetPositions() ([]types.Position, error) {

	var positions []types.Position

	// Get the user profile so we can loop through the accounts.
	userProfile, err := t.GetUserProfile()

	if err != nil {
		return positions, err
	}

	// Get the JSON Per account
	for _, row := range userProfile.Accounts {

		jsonRt, err := t.SendGetRequest("/accounts/" + row.AccountNumber + "/positions")

		if err != nil {
			return positions, err
		}

		// Parse return.
		pos, err := t.parsePositionsJson(jsonRt, row.AccountNumber)

		if err != nil {
			return positions, err
		}

		// Merge back into the master array
		for _, row := range pos {
			positions = append(positions, row)
		}
	}

	return positions, nil
}

//
// Process JSON returned from an position api call.
//
func (t *Api) parsePositionsJson(body string, accountNumber string) ([]types.Position, error) {

	var positions []types.Position

	// Must have some data.
	vo0 := gjson.Get(body, "positions.position")

	// Only one account
	if !vo0.Exists() {
		return positions, nil
	}

	// Do we have only one account?
	vo := gjson.Get(body, "positions.position.id")

	// Only one account
	if vo.Exists() {

		var ws types.Position

		if err := json.Unmarshal([]byte(gjson.Get(body, "positions.position").String()), &ws); err != nil {
			return positions, nil
		}

		// Add in account number
		ws.AccountId = accountNumber

		// Add to master array
		positions = append(positions, ws)

	} else { // More than one accounts

		vo2 := gjson.Get(body, "positions.position")

		// Loop through the different accounts.
		vo2.ForEach(func(key, value gjson.Result) bool {

			var ws types.Position

			if err := json.Unmarshal([]byte(value.String()), &ws); err != nil {
				return false
			}

			// Add in account number
			ws.AccountId = accountNumber

			// Add to master array
			positions = append(positions, ws)

			// Return happy from the loop
			return true
		})

	}

	// Return happy
	return positions, nil
}
