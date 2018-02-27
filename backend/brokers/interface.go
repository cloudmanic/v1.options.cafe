package brokers

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Api interface {
	GetBrokerConfig() *types.BrokerConfig
	SendGetRequest(string) (string, error)
	GetBalances() ([]types.Balance, error)
	GetHistoryByAccountId(string) ([]types.History, error)
	GetOrders() ([]types.Order, error)
	GetAllOrders() ([]types.Order, error)
	GetQuotes([]string) ([]types.Quote, error)
	GetUserProfile() (types.UserProfile, error)
	DoRefreshAccessTokenIfNeeded(models.User) error
	GetTimeSalesQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error)
	GetHistoricalQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error)
}

/* End File */
