//
// Date: 2018-04-03
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package types

type OptionsChain struct {
	Underlying     string             `json:"underlying"`
	ExpirationDate Date               `json:"expiration_date"`
	Puts           []OptionsChainItem `json:"puts"`
	Calls          []OptionsChainItem `json:"calls"`
}

type OptionsChainItem struct {
	Underlying       string  `json:"underlying"`
	Symbol           string  `json:"symbol"`
	OptionType       string  `json:"option_type"`
	Description      string  `json:"description"`
	Strike           float64 `json:"strike"`
	ExpirationDate   Date    `json:"expiration_date"`
	Last             float64 `json:"last"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"change_percentage"`
	Volume           int     `json:"volume"`
	AverageVolume    int     `json:"average_volume"`
	LastVolume       int     `json:"last_volume"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Close            float64 `json:"close"`
	Bid              float64 `json:"bid"`
	BidSize          int     `json:"bid_size"`
	Ask              float64 `json:"ask"`
	AskSize          int     `json:"ask_size""`
	OpenInterest     int     `json:"open_interest"`
	ImpliedVol       float64 `json:"implied_vol""`
	Delta            float64 `json:"delta""`
	Gamma            float64 `json:"gamma""`
	Theta            float64 `json:"theta""`
	Vega             float64 `json:"vega""`
}

/* End File */
