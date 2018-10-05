//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type NotifyChannel struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserId    uint      `gorm:"index" sql:"not null;index:UserId" json:"-"`
	Type      string    `sql:"not null;type:ENUM('One Signal'); default:'One Signal'" json:"type"`
	ChannelId string    `sql:"not null" json:"channel_id"`
}

//
// Validate for this model.
//
func (a NotifyChannel) Validate(db Datastore) error {
	return validation.ValidateStruct(&a,

		// Type
		validation.Field(&a.Type, validation.Required.Error("The type field is required."), validation.By(validateType)),

		// ChannelId
		validation.Field(&a.ChannelId, validation.Required.Error("The channel_id field is required.")),
	)
}

//
// Validate Type
//
func validateType(value interface{}) error {
	if value != "One Signal" {
		return errors.New("Field channel_id must be set to One Signal.")
	}
	return nil
}

/* End File */
