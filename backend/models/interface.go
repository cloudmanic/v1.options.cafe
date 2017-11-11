//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

// Database interface
type Datastore interface {

	// Broker
	UpdateBroker(broker Broker) error
	GetBrokerById(id uint) (Broker, error)
	GetBrokerTypeAndUserId(userId uint, brokerType string) ([]Broker, error)
	CreateNewBroker(name string, user User, accessToken string, refreshToken string, tokenExpirationDate time.Time) (Broker, error)

	// Forgot Password
	GetUserFromToken(token string) (User, error)
	DeleteForgotPasswordByToken(token string) error
	DoResetPassword(user_email string, ip string) error
	GetForgotPasswordStepOneEmailText(name string, email string, url string) string
	GetForgotPasswordStepOneEmailHtml(name string, email string, url string) string

	// Sessions
	GetByAccessToken(accessToken string) (Session, error)
	CreateSession(UserId uint, UserAgent string, LastIpAddress string) (Session, error)

	// Symbols
	GetAllSymbols() []Symbol
	SearchSymbols(query string) ([]Symbol, error)
	CreateNewSymbol(short string, name string) (Symbol, error)
	UpdateSymbol(id uint, short string, name string) (Symbol, error)

	// Users
	GetAllUsers() []User
	GetAllActiveUsers() []User
	UpdateUser(user *User) error
	GetUserById(id uint) (User, error)
	ValidatePassword(password string) error
	GetUserByEmail(email string) (User, error)
	ResetUserPassword(id uint, password string) error
	ValidateUserLogin(email string, password string) error
	GetUserByStripeCustomer(customerId string) (User, error)
	ValidateCreateUser(first string, last string, email string, password string) error
	LoginUserByEmailPass(email string, password string, userAgent string, ipAddress string) (User, error)
	CreateUser(first string, last string, email string, password string, userAgent string, ipAddress string) (User, error)

	// Watchlists
	GetWatchlistsById(id uint) (Watchlist, error)
	GetWatchlistsByUserId(userId uint) ([]Watchlist, error)
	CreateWatchlist(userId uint, name string) (Watchlist, error)

	CreateNewWatchlist(user User, name string) (Watchlist, error) // TODO: kill this one

	// WatchlistSymbols
	CreateNewWatchlistSymbol(wList Watchlist, symb Symbol, user User, order uint) (WatchlistSymbol, error)
}

/* End File */
