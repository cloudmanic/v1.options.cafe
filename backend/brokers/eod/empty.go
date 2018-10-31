//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: Just placeholder for functions we do not care about.
//
package eod

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

func (t *Api) GetBrokerConfig() *types.BrokerConfig {
	return nil
}

func (t *Api) GetBalances() ([]types.Balance, error) {
	return []types.Balance{}, nil
}

func (t *Api) GetAllHistory() ([]types.History, error) {
	return []types.History{}, nil
}

func (t *Api) GetHistoryByAccountId(string) ([]types.History, error) {
	return []types.History{}, nil
}

func (t *Api) CancelOrder(accountId string, orderId string) error {
	return nil
}

func (t *Api) SubmitOrder(accountId string, order types.Order) (types.OrderSubmit, error) {
	return types.OrderSubmit{}, nil
}

func (t *Api) PreviewOrder(accountId string, order types.Order) (types.OrderPreview, error) {
	return types.OrderPreview{}, nil
}

func (t *Api) GetOrders() ([]types.Order, error) {
	return []types.Order{}, nil
}

func (t *Api) GetPositions() ([]types.Position, error) {
	return []types.Position{}, nil
}

func (t *Api) GetAllOrders() ([]types.Order, error) {
	return []types.Order{}, nil
}

func (t *Api) GetQuotes([]string) ([]types.Quote, error) {
	return []types.Quote{}, nil
}

func (t *Api) GetUserProfile() (types.UserProfile, error) {
	return types.UserProfile{}, nil
}

func (t *Api) DoRefreshAccessTokenIfNeeded(models.User) error {
	return nil
}

func (t *Api) GetTimeSalesQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error) {
	return []types.HistoryQuote{}, nil
}

func (t *Api) GetHistoricalQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error) {
	return []types.HistoryQuote{}, nil
}

/* End File */
