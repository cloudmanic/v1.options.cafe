package brokers

import (
	"app.options.cafe/backend/brokers/types"
	"app.options.cafe/backend/models"
)

type Api interface {
	SetActiveSymbols([]string)
	SendGetRequest(string) (string, error)
	GetBalances() ([]types.Balance, error)
	GetHistoryByAccountId(string) ([]types.History, error)
	GetMarketStatus() (types.MarketStatus, error)
	GetOrders() ([]types.Order, error)
	GetAllOrders() ([]types.Order, error)
	GetQuotes([]string) ([]types.Quote, error)
	GetUserProfile() (types.UserProfile, error)
	DoRefreshAccessTokenIfNeeded(models.User) error
}

/* End File */
