//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package types

type Position struct {
	Id        uint    `json:"id"`
	AccountId string  `json:"account_id"`
	OpenDate  string  `json:"date_acquired"`
	Quantity  float64 `json:"quantity"`
	CostBasis float64 `json:"cost_basis"`
	Symbol    string  `json:"symbol"`
}

/* End File */
