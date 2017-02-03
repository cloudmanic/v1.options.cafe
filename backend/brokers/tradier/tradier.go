package tradier

import (
  "sync"
  "strings"  
)

const (
	apiBaseUrl = "https://api.tradier.com/v1"
)

type Api struct {
  mu sync.Mutex
  defaultAccountId string
  
  activeSymbols string
  ApiKey string
  oldActiveSymbols string
}

//
// Add symbols to active list.
//
func (t * Api) SetActiveSymbols(symbols []string) {
  
  // Lock da memory
	t.mu.Lock()
	defer t.mu.Unlock()    
  
  t.activeSymbols = strings.Join(symbols, ",")
  
}

//
// Set the default account.
//
func (t * Api) SetDefaultAccountId(accountId string) {
  
  // Lock da memory
	t.mu.Lock()
	defer t.mu.Unlock()    
	
  // Store value
  t.defaultAccountId = accountId
    
}

/* End File */