package brokers

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Api interface {
	GetBrokerConfig() *types.BrokerConfig
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
	GetHistoricalQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error)
}

/* End File */
