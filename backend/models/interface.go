//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	"app.options.cafe/brokers/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

// Database interface
type Datastore interface {

	// Gorm Functions
	New() *gorm.DB

	// Applications
	ValidateClientIdGrantType(clientId string, grantType string) (Application, error)

	// Generic database functions
	Count(model interface{}, params QueryParam) (uint, error)
	Query(model interface{}, params QueryParam) error
	CreateNewRecord(model interface{}, params InsertParam) error
	QueryWithNoFilterCount(model interface{}, params QueryParam) (int, error)
	GetQueryMetaData(limitCount int, noLimitCount int, params QueryParam) QueryMetaData

	// Settings
	SettingsGetOrCreateByUserId(userId uint) Settings

	// Broker
	UpdateBroker(broker Broker) error
	GetBrokerById(id uint) (Broker, error)
	KickStartBroker(user User, broker Broker)
	GetBrokerTypeAndUserId(userId uint, brokerType string) ([]Broker, error)
	CreateNewBroker(name string, user User, accessToken string, refreshToken string, tokenExpirationDate time.Time) (Broker, error)

	// Broker Accounts
	GetBrokerAccountByIdUserId(id uint, userId uint) (BrokerAccount, error)
	GetBrokerFromBrokerAccountAndUserId(id uint, userId uint) (Broker, error)
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
	CreateSession(UserId uint, appId uint, UserAgent string, LastIpAddress string) (Session, error)

	// Symbols
	GetAllSymbols() []Symbol
	GetSymbolByShortName(short string) (Symbol, error)
	ValidateSymbolId(value interface{}) error
	SearchSymbols(query string, sType string) ([]Symbol, error)
	CreateNewOptionSymbol(short string) (Symbol, error)
	CreateNewSymbol(short string, name string, sType string) (Symbol, error)
	UpdateSymbol(id uint, short string, name string, sType string) (Symbol, error)
	LoadSymbolsByOptionsChain(chain types.OptionsChain) error
	GetOptionByParts(optionUnderlying string, optionType string, optionExpire time.Time, optionStrike float64) (Symbol, error)

	// Users
	GetAllUsers() []User
	GetAllActiveUsers() []User
	GetAllActiveOrTrialUsers() []User
	DeleteUser(user *User) error
	UpdateUser(user *User) error
	GetUserById(id uint) (User, error)
	ValidatePassword(password string) error
	GetUserByEmail(email string) (User, error)
	GetUserByGoogleSubId(sub string) (User, error)
	SetDefaultWatchList(user User)
	ResetUserPassword(id uint, password string) error
	ValidateUserLogin(email string, password string) error
	GetUserByStripeCustomer(customerId string) (User, error)
	DeleteUserWithStripe(user User) error
	ValidateUserEmail(userId uint, email string) error
	ApplyCoupon(user User, couponCode string) error
	UpdateCreditCard(user User, stripeToken string) error
	GetInvoiceHistoryWithStripe(user User) ([]UserInvoice, error)
	GetSubscriptionWithStripe(user User) (UserSubscription, error)
	CreateNewUserWithStripe(user User, plan string, token string, coupon string) error
	ValidateCreateUser(first string, last string, email string, googleAuth bool) error
	LoginUserById(id uint, appId uint, userAgent string, ipAddress string) (User, error)
	LoginUserByEmailPass(email string, password string, appId uint, userAgent string, ipAddress string) (User, error)
	CreateUser(first string, last string, email string, password string, appId uint, userAgent string, ipAddress string) (User, error)
	CreateUserFromGoogle(first string, last string, email string, subId string, appId uint, userAgent string, ipAddress string) (User, error)

	// Orders
	GetOrdersByUser(userId uint) []Order
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
	WatchlistDeleteById(id uint) error
	WatchlistReorder(id uint, ids []int) error
	WatchlistUpdate(id uint, name string) error
	GetWatchlistsById(id uint) (Watchlist, error)
	WatchlistRemoveSymbol(id uint, symbId uint) error
	GetWatchlistsByUserId(userId uint) ([]Watchlist, error)
	CreateWatchlist(userId uint, name string) (Watchlist, error)
	GetWatchlistsByIdAndUserId(id uint, userId uint) (Watchlist, error)

	// WatchlistSymbols
	PrependWatchlistSymbol(w *WatchlistSymbol) error
	WatchlistSymbolGetByWatchlistId(id uint) []WatchlistSymbol
	CreateWatchlistSymbol(wList Watchlist, symb Symbol, user User, order uint) (WatchlistSymbol, error)

	// ActiveSymbol
	GetActiveSymbolsByUser(userId uint) ([]ActiveSymbol, error)
	CreateActiveSymbol(userId uint, symbol string) (ActiveSymbol, error)

	// Screener
	GetScreenersByUserId(userId uint) ([]Screener, error)
	GetScreenerByIdAndUserId(id uint, userId uint) (Screener, error)
}

/* End File */
