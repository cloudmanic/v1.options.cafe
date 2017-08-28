//
// Date: 8/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
  "os"
  "net/url"
  "net/http"
)

//
// Subscribe to a sendy newsletter list
//
func SendySubscribe(listId string, email string, first string, last string) {
  
  var listIdString = ""
  
  // Get the proper list id from our configs
  switch listId {
    
    case "subscribers":
      if len(os.Getenv("SENDY_SUBSCRIBE_LIST")) > 0 {
        
        listIdString = os.Getenv("SENDY_SUBSCRIBE_LIST") 
        
      }

    case "no-brokers":
      if len(os.Getenv("SENDY_NO_BROKER_LIST")) > 0 {
        
        listIdString = os.Getenv("SENDY_NO_BROKER_LIST") 
        
      }    
    
  }
  
  // Make sure we have a list id.
  if len(listIdString) == 0 {
    MajorLog("No listIdString found in SendySubscribe : " + listId + " - " + listIdString)
    return 
  }
  
  // Build form request
  form := url.Values{
    "list": {listIdString},
    "email": {email},
    "name": {first + " " + last},
    "FirstName": {first},
    "LastName": {last},		
  }
  
  // Send request.
  resp, err := http.PostForm("https://sendy.cloudmanic.com/subscribe", form)
  
  if err != nil {
    Error(err, "SendySubscribe - Unable to subscribe " + email + " to Sendy Subscriber list.")
  }
  
  if resp.StatusCode != http.StatusOK {
    Error(err, "SendySubscribe (no 200) - Unable to subscribe " + email + " to Sendy Subscriber list.")    
  }
  
  defer resp.Body.Close()   
  
}

/* End File */