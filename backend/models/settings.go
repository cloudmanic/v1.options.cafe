//
// Date: 7/22/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Settings struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserId    uint      `sql:"not null;index:UserId" json:"user_id"`

	// Put credit Spread
	StrategyPcsClosePrice float64 `json:"strategy_pcs_close_price" sql:"not null;default:0.03"`
	StrategyPcsOpenPrice  string  `json:"strategy_pcs_open_price" sql:"not null;default:'mid-point'"`
	StrategyPcsLots       uint    `json:"strategy_pcs_lots" sql:"not null;default:10"`

	// Call credit Spread
	StrategyCcsClosePrice float64 `json:"strategy_ccs_close_price" sql:"not null;default:0.03"`
	StrategyCcsOpenPrice  string  `json:"strategy_ccs_open_price" sql:"not null;default:'mid-point'"`
	StrategyCcsLots       uint    `json:"strategy_ccs_lots" sql:"not null;default:10"`

	// Put debit Spread
	StrategyPdsClosePrice float64 `json:"strategy_pds_close_price" sql:"not null;default:0.03"`
	StrategyPdsOpenPrice  string  `json:"strategy_pds_open_price" sql:"not null;default:'mid-point'"`
	StrategyPdsLots       uint    `json:"strategy_pds_lots" sql:"not null;default:10"`

	// Call debit Spread
	StrategyCdsClosePrice float64 `json:"strategy_cds_close_price" sql:"not null;default:0.03"`
	StrategyCdsOpenPrice  string  `json:"strategy_cds_open_price" sql:"not null;default:'mid-point'"`
	StrategyCdsLots       uint    `json:"strategy_cds_lots" sql:"not null;default:10"`

	// Trade Filled
	NoticeTradeFilledEmail string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_trade_filled_email"`
	NoticeTradeFilledSms   string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_trade_filled_sms"`
	NoticeTradeFilledPush  string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_trade_filled_push"`

	// Market Open
	NoticeMarketOpenedEmail string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_open_email"`
	NoticeMarketOpenedSms   string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_open_sms"`
	NoticeMarketOpenedPush  string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_open_push"`

	// Market Closed
	NoticeMarketClosedEmail string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_closed_email"`
	NoticeMarketClosedSms   string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_closed_sms"`
	NoticeMarketClosedPush  string `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"notice_market_closed_push"`
}

//
// Validate for this model.
//
func (a Settings) Validate(db Datastore, userId uint) error {

	// Return validation
	return validation.ValidateStruct(&a,

		// Strategies - Vertical Spreads
		validation.Field(&a.StrategyPcsClosePrice, validation.Required.Error("The strategy_pcs_close_price field is required.")),
		validation.Field(&a.StrategyPcsLots, validation.Required.Error("The strategy_pcs_close_price field is required.")),

		validation.Field(&a.StrategyCcsClosePrice, validation.Required.Error("The strategy_ccs_close_price field is required.")),
		validation.Field(&a.StrategyCcsLots, validation.Required.Error("The strategy_ccs_lots field is required.")),

		validation.Field(&a.StrategyPdsClosePrice, validation.Required.Error("The strategy_pds_close_price field is required.")),
		validation.Field(&a.StrategyPdsLots, validation.Required.Error("The strategy_pds_lots field is required.")),

		validation.Field(&a.StrategyCdsClosePrice, validation.Required.Error("The strategy_cds_close_price field is required.")),
		validation.Field(&a.StrategyCdsLots, validation.Required.Error("The strategy_cds_lots field is required.")),

		validation.Field(&a.StrategyPcsOpenPrice,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("mid-point", "bid", "ask").Error("The strategy_pcs_open_price must be mid-point, bid, or ask."),
		),

		validation.Field(&a.StrategyCcsOpenPrice,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("mid-point", "bid", "ask").Error("The strategy_ccs_open_price must be mid-point, bid, or ask."),
		),

		validation.Field(&a.StrategyPdsOpenPrice,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("mid-point", "bid", "ask").Error("The strategy_pds_open_price must be mid-point, bid, or ask."),
		),

		validation.Field(&a.StrategyCdsOpenPrice,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("mid-point", "bid", "ask").Error("The strategy_cds_open_price must be mid-point, bid, or ask."),
		),

		// NoticeTradeFilledEmail
		validation.Field(&a.NoticeTradeFilledEmail,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_trade_filled_email must be Yes or No."),
		),

		// NoticeTradeFilledSms
		validation.Field(&a.NoticeTradeFilledSms,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_open_sms must be Yes or No."),
		),

		// NoticeTradeFilledPush
		validation.Field(&a.NoticeTradeFilledPush,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_open_push must be Yes or No."),
		),

		// NoticeMarketOpenedEmail
		validation.Field(&a.NoticeMarketOpenedEmail,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_trade_filled_email must be Yes or No."),
		),

		// NoticeMarketOpenedSms
		validation.Field(&a.NoticeMarketOpenedSms,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_open_sms must be Yes or No."),
		),

		// NoticeMarketOpenedPush
		validation.Field(&a.NoticeMarketOpenedPush,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_open_push must be Yes or No."),
		),

		// NoticeMarketClosedEmail
		validation.Field(&a.NoticeMarketClosedEmail,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_closed_email must be Yes or No."),
		),

		// NoticeMarketClosedSms
		validation.Field(&a.NoticeMarketClosedSms,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_closed_sms must be Yes or No."),
		),

		// NoticeMarketClosedPush
		validation.Field(&a.NoticeMarketClosedPush,
			validation.Required.Error("The notice_trade_filled_email field is required."),
			validation.In("Yes", "No").Error("The notice_market_closed_push must be Yes or No."),
		),
	)
}

//
// Get or create by users id.
//
func (t *DB) SettingsGetOrCreateByUserId(userId uint) Settings {

	var s Settings

	// First make sure we don't already have this settings
	if t.Where("user_id = ?", userId).First(&s).RecordNotFound() {

		// Create entry.
		s = Settings{UserId: userId}

		t.Create(&s)

	}

	// Return the settings.
	return s
}

/* End File */
