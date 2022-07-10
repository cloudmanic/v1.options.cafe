//
// Date: 7/19/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"encoding/json"

	"app.options.cafe/brokers/types"
	"github.com/tidwall/gjson"
)

//
// GetGainsByAccountId will return gains and losses by Account Id
//
func (t *Api) GetGainsByAccountId(accountId string) ([]types.Gain, error) {
	var gains []types.Gain

	// Get the JSON - Set big limit to get all data
	jsonRt, err := t.SendGetRequest("/accounts/" + accountId + "/gainloss?limit=10")

	if err != nil {
		return gains, err
	}

	// Make sure we have at least one gain
	vo := gjson.Get(jsonRt, "gainloss.closed_position")

	if !vo.Exists() {
		return gains, nil
	}

	// Do we have only one gain or more?
	vo2 := gjson.Get(jsonRt, "gainloss.closed_position.close_date")

	// Only one event
	if vo2.Exists() {
		var gain types.Gain

		if err := json.Unmarshal([]byte(vo.String()), &gain); err != nil {
			return gains, err
		}

		gains = append(gains, gain)
	} else {
		if err := json.Unmarshal([]byte(vo.String()), &gains); err != nil {
			return gains, err
		}
	}

	// Return happy
	return gains, nil
}
