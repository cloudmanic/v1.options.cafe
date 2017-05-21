package brokers

import (
  "app.options.cafe/backend/brokers/types"    
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
  GetWatchLists() ([]types.Watchlist, error)
  GetWatchList(string) (types.Watchlist, error)
  
}

/* End File */