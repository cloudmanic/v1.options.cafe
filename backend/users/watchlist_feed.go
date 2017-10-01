//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
  "encoding/json"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/controllers" 
  "app.options.cafe/backend/library/services"
)

//
// Send a user's watchlist up the websocket channel
//
func WsSendWatchlists(user *User) {
  
  // Get the watchlists
  wLists, err := DB.GetWatchlistsByUserId(user.Profile.Id)
  
  if err != nil {
    return
  }  
  
  type Wl struct {
    Id uint    
    Name string
    List []string
  }  
  
  // Loop through the different watchlists
  for _, row := range wLists {
    
    // Clean up data we send.
    l := &Wl{ Id: row.Id, Name: row.Name }
    
    for _, row2 := range row.Symbols {
      l.List = append(l.List, row2.Symbol.ShortName)
    }

    // Convert to a json string.
    dataJson, err := json.Marshal(l)

    if err != nil {
      services.Error(err, "WsSendWatchlists() json.Marshal (#1)")
      continue
    } 
    
    // Build JSON we send
    jsonSend, err := WsSendJsonBuild("Watchlist:refresh", dataJson)    
        
    if err != nil {
      services.Error(err, "WsSendWatchlists() WsSendJsonBuild (#2)")
      continue
    }    
    
    // Send up the websocket
    user.DataChan <- controllers.SendStruct{ UserId: user.Profile.Id, Message: string(jsonSend) }
  
  }
 
  // Return happy
  return 
}

//
// Verify we have default watchlist in place.
//
func VerifyDefaultWatchList(user models.User) {
  
  // Setup defaults.
  type Y struct {
    SymShort string
    SymLong string
  }
  
  var m []Y
  m = append(m, Y{ SymShort: "SPY", SymLong: "SPDR S&P 500" })
  m = append(m, Y{ SymShort: "IWM", SymLong: "Ishares Russell 2000 Etf" })
  m = append(m, Y{ SymShort: "VIX", SymLong: "CBOE Volatility S&P 500 Index" })
  m = append(m, Y{ SymShort: "AMZN", SymLong: "Amazon.com Inc" })
  m = append(m, Y{ SymShort: "AAPL", SymLong: "Apple Inc." })      
  m = append(m, Y{ SymShort: "SBUX", SymLong: "Starbucks Corp" })
  m = append(m, Y{ SymShort: "BAC", SymLong: "Bank Of America Corporation" })

  // See if this user already had a watchlist
  _, err := DB.GetWatchlistsByUserId(user.Id)
  
  // If no watchlists we create a default one with some default symbols.  
  if err != nil {

    wList, err := DB.CreateNewWatchlist(user, "Default")

    if err != nil {
      services.Error(err, "(CreateNewWatchlist) Unable to create watchlist Default")
      return
    }

    for key, row := range m {

      // Add some default symbols - SPY
      symb, err := DB.CreateNewSymbol(row.SymShort, row.SymLong)
      
      if err != nil {
        services.Error(err, "(VerifyDefaultWatchList) Unable to create symbol " + row.SymShort)
        return
      }
      
      // Add lookup
      _, err2 := DB.CreateNewWatchlistSymbol(wList, symb, user, uint(key))      
  
      if err2 != nil {
        services.Error(err2, "(CreateNewWatchlistSymbol) Unable to create symbol " + row.SymShort + " lookup")
        return
      }
    
    }
    
  }
  
  return
  
}

/* End File */