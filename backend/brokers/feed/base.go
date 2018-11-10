package feed

import (
	"sync"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
)

type Base struct {
	User        models.User
	Api         brokers.Api
	DB          models.Datastore
	BrokerId    uint
	WsWriteChan chan websocket.SendStruct
	Polling     bool
	MuPolling   sync.Mutex
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

	_, err := t.Api.GetPositions()

	if err != nil {
		services.Warning(err)
	}

	// Setup tickers for broker polling.
	go t.DoOrdersTicker()
	go t.DoGetHistoryTicker()
	go t.DoUserProfileTicker()
	go t.DoGetBalancesTicker()
	go t.DoAccessTokenRefresh()
}

/* End File */
