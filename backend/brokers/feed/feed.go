package feed

import (
	"sync"
	"time"

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

	muPositions sync.Mutex
	Positions   []types.Position

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
	go t.DoPositionsTicker()
	go t.DoUserProfileTicker()
	go t.DoGetDetailedQuotes()
	go t.DoGetMarketStatusTicker()
	go t.DoGetBalancesTicker()
	go t.DoAccessTokenRefresh()

	// Do Archive Calls
	//go t.DoOrdersArchive()

}

// ---------------------- Tickers (polling) ---------------------------- //

//
// Ticker - Positions : 3 seconds
//
func (t *Base) DoPositionsTicker() {

	var err error

	for {

		// Load up positions
		err = t.GetPositions()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 3 second.
		time.Sleep(time.Second * 3)

	}

}

//
// Ticker - Get DetailedQuotes : 1 second
//
func (t *Base) DoGetDetailedQuotes() {

	for {

		// Load up our DetailedQuotes
		err := t.GetActiveSymbolsDetailedQuotes()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 1 second
		time.Sleep(time.Second * 1)

	}

}

//
// Ticker - Get GetMarketStatus : 5 seconds
//
func (t *Base) DoGetMarketStatusTicker() {

	var err error

	for {

		// Load up market status.
		err = t.GetMarketStatus()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 10 second.
		time.Sleep(time.Second * 5)

	}

}

// Ticker - Get GetBalances : 5 seconds
//
func (t *Base) DoGetBalancesTicker() {

	var err error

	for {

		// Load up market status.
		err = t.GetBalances()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 5 second.
		time.Sleep(time.Second * 5)

	}

}

//
// Ticker - See if we need to refresh an access token : 60 seconds
//
func (t *Base) DoAccessTokenRefresh() {

	var err error

	for {

		err = t.AccessTokenRefresh()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 60 second.
		time.Sleep(time.Second * 60)

	}

}

/* End File */
