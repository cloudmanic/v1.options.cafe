package types

type Watchlist struct {
  Name string 
  Id string
  Symbols []WatchlistSymbol      
}

type WatchlistSymbol struct {
  Id string `json:"id"`
  Name string `json:"symbol"`
}