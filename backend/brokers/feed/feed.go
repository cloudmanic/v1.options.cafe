package feed

import (
	"sync"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/controllers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Base struct {
	User     models.User
	Api      brokers.Api
	DB       models.Datastore
	BrokerId uint

	DataChan  chan controllers.SendStruct
	QuoteChan chan controllers.SendStruct

	muOrders sync.Mutex
	Orders   []types.Order

	muBalances sync.Mutex
	Balances   []types.Balance

	muMarketStatus sync.Mutex
	MarketStatus   types.MarketStatus

	muUserProfile sync.Mutex
	UserProfile   types.UserProfile
}

type SendStruct struct {
	Uri  string `json:"uri"`
	Body string `json:"body"`
}

//
// When we have a broker access token and an active user we call this.
// We start fetching data from the broker and such. This continues to run
// until the broker access token stops working, or the user expires
// or is revoked.
//
func (t *Base) Start() {

	services.Info("Starting Polling....")

	// Setup tickers for broker polling.
	go t.DoOrdersTicker()
	go t.DoUserProfileTicker()
	go t.DoGetDetailedQuotes()
	go t.DoGetMarketStatusTicker()
	go t.DoGetBalancesTicker()
	go t.DoAccessTokenRefresh()
	go t.DoOrdersArchive()
}

/* End File */
