package tradier

import (
  "sync"
  "strings"  
)

const (
	apiBaseUrl = "https://api.tradier.com/v1"
)

type Api struct {
  muActiveSymbols sync.Mutex
  activeSymbols string
  
  ApiKey string
}

//
// Add symbols to active list.
//
func (t * Api) SetActiveSymbols(symbols []string) {
  
  // Lock da memory
	t.muActiveSymbols.Lock()
	defer t.muActiveSymbols.Unlock()    
  
  t.activeSymbols = strings.Join(symbols, ",")
  
}

/* End File */