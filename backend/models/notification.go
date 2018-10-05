//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Notification struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	UserId       uint      `gorm:"index" sql:"not null;index:UserId" json:"-"`
	Status       string    `sql:"not null;type:ENUM('Pending', 'Sent', 'Seen'); default:'Pending'" json:"type"`
	Channel      string    `sql:"not null;type:ENUM('Push', 'SMS', 'Email', 'Slack', 'Zapier'); default:'Push'" json:"channel"`
	Object       string    `sql:"not null" json:"object"`
	ObjectId     uint      `sql:"not null" json:"object_id"`
	Title        string    `sql:"not null" json:"title"`
	ShortMessage string    `sql:"not null" json:"short_message"`
	LongMessage  string    `sql:"not null" json:"long_message"`
	Expires      time.Time `json:"-"`
}

/* End File */
