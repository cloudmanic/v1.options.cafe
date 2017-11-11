package models

import (
	"errors"
	"time"

	"app.options.cafe/backend/library/services"
)

type Session struct {
	Id            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserId        uint   `sql:"not null;index:UserId"`
	UserAgent     string `sql:"not null"`
	AccessToken   string `sql:"not null"`
	LastIpAddress string `sql:"not null"`
	LastActivity  time.Time
}

//
// Get by Access token.
//
func (t *DB) GetByAccessToken(accessToken string) (Session, error) {

	session := Session{}

	if t.First(&session, "access_token = ?", accessToken).RecordNotFound() {
		return Session{}, errors.New("Access Token Not Found - Unable to Authenticate")
	}

	// Double check because of case sensitivity
	if session.AccessToken == accessToken {
		return Session{}, errors.New("Access Token Not Found - Unable to Authenticate")
	}

	// Return happy
	return session, nil
}

//
// Create new session. A user can have more than one session. Typically it is one session per browser or device.
// We return the session object. The big thing here is we create the access token for this session.
//
func (t *DB) CreateSession(UserId uint, UserAgent string, LastIpAddress string) (Session, error) {

	// Create an access token.
	access_token, err := t.GenerateRandomString(50)

	if err != nil {
		services.Error(err, "CreateSession - Unable to create random string (access_token)")
		return Session{}, err
	}

	// Save the session into the database.
	session := Session{UserId: UserId, UserAgent: UserAgent, AccessToken: access_token, LastIpAddress: LastIpAddress}
	t.Create(&session)

	// Return the session.
	return session, nil
}

/* End File */
