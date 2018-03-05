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

	// Generic database functions
	Count(model interface{}, params QueryParam) (uint, error)
	Query(model interface{}, params QueryParam) error
	QueryWithNoFilterCount(model interface{}, params QueryParam) (int, error)
	GetQueryMetaData(limitCount int, noLimitCount int, params QueryParam) QueryMetaData

	// Broker
	UpdateBroker(broker Broker) error
	GetBrokerById(id uint) (Broker, error)
	GetBrokerTypeAndUserId(userId uint, brokerType string) ([]Broker, error)
	CreateNewBroker(name string, user User, accessToken string, refreshToken string, tokenExpirationDate time.Time) (Broker, error)

	// Broker Accounts
	UpdateBrokerAccount(brokerAccount *BrokerAccount) error
	FirstOrCreateBrokerAccount(brokerAccount *BrokerAccount) (bool, error)
	GetBrokerAccountByBrokerAccountNumber(brokerId uint, accountNumber string) (BrokerAccount, error)

	// Forgot Password
	GetUserFromToken(token string) (User, error)
	DeleteForgotPasswordByToken(token string) error
	DoResetPassword(user_email string, ip string) error
	GetForgotPasswordStepOneEmailText(name string, email string, url string) string
	GetForgotPasswordStepOneEmailHtml(name string, email string, url string) string

	// Sessions
	UpdateSession(session *Session) error
	GetByAccessToken(accessToken string) (Session, error)
	CreateSession(UserId uint, UserAgent string, LastIpAddress string) (Session, error)

	// Symbols
	GetAllSymbols() []Symbol
	SearchSymbols(query string, sType string) ([]Symbol, error)
	CreateNewOptionSymbol(short string) (Symbol, error)
	CreateNewSymbol(short string, name string, sType string) (Symbol, error)
	UpdateSymbol(id uint, short string, name string, sType string) (Symbol, error)

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

	// Orders
	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	HasOrderByBrokerRefUserId(brokerRef string, userId uint) bool
	GetOrdersByUserClassStatusReviewed(userId uint, class string, status string, reviewed string) ([]Order, error)

	// Positions
	CreatePosition(position *Position) error
	UpdatePosition(position *Position) error
	GetPositionByUserSymbolStatusAccount(userId uint, symbolId uint, status string, brokerAccountId uint) (Position, error)

	// TradeGroup
	GetTradeGroups(params QueryParam) ([]TradeGroup, QueryMetaData, error)
	GetTradeGroupById(id uint) (TradeGroup, error)
	CreateTradeGroup(tg *TradeGroup) error
	UpdateTradeGroup(tg *TradeGroup) error

	// Watchlists
	GetWatchlistsById(id uint) (Watchlist, error)
	GetWatchlistsByUserId(userId uint) ([]Watchlist, error)
	CreateWatchlist(userId uint, name string) (Watchlist, error)
	GetWatchlistsByIdAndUserId(id uint, userId uint) (Watchlist, error)

	CreateNewWatchlist(user User, name string) (Watchlist, error) // TODO: kill this one

	// WatchlistSymbols
	CreateNewWatchlistSymbol(wList Watchlist, symb Symbol, user User, order uint) (WatchlistSymbol, error)

	// ActiveSymbol
	GetActiveSymbolsByUser(userId uint) ([]ActiveSymbol, error)
	CreateActiveSymbol(userId uint, symbol string) (ActiveSymbol, error)
}

/* End File */
