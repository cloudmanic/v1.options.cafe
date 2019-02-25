//
// Date: 2018-04-03
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package types

type OptionsChain struct {
	Id             uint               `gorm:"primary_key" json:"_"`
	Underlying     string             `sql:"not null" json:"underlying"`
	UnderlyingLast float64            `sql:"not null" json:"lunderlying_last"`
	ExpirationDate Date               `gorm:"type:date" sql:"not null" json:"expiration_date"`
	Puts           []OptionsChainItem `json:"puts"`
	Calls          []OptionsChainItem `json:"calls"`
}

type OptionsChainItem struct {
	OptionsChainId   uint    `sql:"not null;index:OptionsChainId" json:"_"`
	Underlying       string  `sql:"not null" json:"underlying"`
	Symbol           string  `sql:"not null" json:"symbol"`
	OptionType       string  `sql:"not null" json:"option_type"`
	Description      string  `sql:"not null" json:"description"`
	Strike           float64 `sql:"not null" json:"strike"`
	ExpirationDate   Date    `sql:"not null" json:"expiration_date"`
	Last             float64 `sql:"not null" json:"last"`
	Change           float64 `sql:"not null" json:"change"`
	ChangePercentage float64 `sql:"not null" json:"change_percentage"`
	Volume           int     `sql:"not null" json:"volume"`
	AverageVolume    int     `sql:"not null" json:"average_volume"`
	LastVolume       int     `sql:"not null" json:"last_volume"`
	Open             float64 `sql:"not null" json:"open"`
	High             float64 `sql:"not null" json:"high"`
	Low              float64 `sql:"not null" json:"low"`
	Close            float64 `sql:"not null" json:"close"`
	Bid              float64 `sql:"not null" json:"bid"`
	BidSize          int     `sql:"not null" json:"bid_size"`
	Ask              float64 `sql:"not null" json:"ask"`
	AskSize          int     `sql:"not null" json:"ask_size""`
	OpenInterest     int     `sql:"not null" json:"open_interest"`
	ImpliedVol       float64 `sql:"not null" json:"implied_vol""`
	Delta            float64 `sql:"not null" json:"delta""`
	Gamma            float64 `sql:"not null" json:"gamma""`
	Theta            float64 `sql:"not null" json:"theta""`
	Vega             float64 `sql:"not null" json:"vega""`
}

/* End File */
