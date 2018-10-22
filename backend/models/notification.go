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
	Status       string    `sql:"not null;type:ENUM('pending', 'sent', 'seen', 'expired'); default:'pending'" json:"status"`
	Channel      string    `sql:"not null;type:ENUM('web-push', 'sms-push', 'email', 'slack', 'zapier', 'in-app'); default:'web-push'" json:"channel"`
	Uri          string    `sql:"not null" json:"uri"`
	UriRefId     uint      `sql:"not null" json:"uri_ref_id"`
	Title        string    `sql:"not null" json:"title"`
	ShortMessage string    `sql:"not null" json:"short_message"`
	LongMessage  string    `sql:"type:text" json:"long_message"`
	SentTime     time.Time `json:"sent_time"`
	Expires      time.Time `json:"-"`
}

/* End File */
