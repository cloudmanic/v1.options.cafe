package brokers

import (
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
)

type Api interface {
	GetBrokerConfig() *types.BrokerConfig
	GetBalances() ([]types.Balance, error)
	GetAllHistory() ([]types.History, error)
	GetHistoryByAccountId(string) ([]types.History, error)
	CancelOrder(accountId string, orderId string) error
	SubmitOrder(accountId string, order types.Order) (types.OrderSubmit, error)
	PreviewOrder(accountId string, order types.Order) (types.OrderPreview, error)
	GetOrders() ([]types.Order, error)
	GetPositions() ([]types.Position, error)
	GetAllOrders() ([]types.Order, error)
	GetQuotes([]string) ([]types.Quote, error)
	GetUserProfile() (types.UserProfile, error)
	DoRefreshAccessTokenIfNeeded(models.User) error
	GetOptionsExpirationsBySymbol(symb string) ([]string, error)
	GetOptionsChainByExpiration(symbol string, expireDate string) (types.OptionsChain, error)
	GetTimeSalesQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error)
	GetHistoricalQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error)
}

/* End File */
