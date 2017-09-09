//
// Date: 9/6/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
  "os"
  "time"
  "errors"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/library/helpers"
  "app.options.cafe/backend/library/services"   
)

var (
  genericError = []byte("Something went wrong while authorizing your account. Please try again or contact help@options.cafe. Sorry for the trouble.")
)

// Json response object
type tokenResponse struct {
  Token string `json:"access_token"`
  ExpiresSec int64  `json:"expires_in"`
  IssueDateStr string `json:"issued_at"`
  RefreshToken string `json:"refresh_token"`
  Scope string
  Status string
}

//
// Obtain an Authorization Code - http://localhost:7652/tradier/authorize?user=1
//
func DoAuthCode(w http.ResponseWriter, r *http.Request) {
  
  // Make sure we have a user id.
  userId := r.URL.Query().Get("user")
  
  if userId == "" {
    var msg = "Tradier - DoAuthCode - No user id provided."
    services.Error(errors.New(msg), msg)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return 
  }  
  
  // Start the database
  var DB = models.DB{}
  DB.Start()
  defer DB.Connection.Close()  
  
  // Make sure this is a valid user.
  u, _ := strconv.ParseUint(userId, 10, 32)
  user, err := DB.GetUserById(uint(u))
  
  if err != nil {
    services.Error(err, "Tradier - DoAuthCode - No user found.")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return
  }
  
  // Log
  services.Log("Tradier authorization starting for " + user.Email)
    
  // Redirect to tradier to auth
  var url = apiBaseUrl + "/oauth/authorize?client_id=" + os.Getenv("TRADIER_CONSUMER_KEY") + "&scope=read,write,market,trade,stream&state=" + strconv.Itoa(int((user.Id)))  
  http.Redirect(w, r, url, 302)
  
}

//
// Do Obtain an Authorization Code Callback - http://localhost:7652/tradier/callback
//
func DoAuthCallback(w http.ResponseWriter, r *http.Request) {
  
  // Make sure we have a code.
  code := r.URL.Query().Get("code")
  
  if code == "" {
    var msg = "Tradier - DoAuthCallback - No auth code provided. (#1)"
    services.Error(errors.New(msg), msg)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return 
  }

  // Make sure we have a state.
  state := r.URL.Query().Get("state")
  
  if state == "" {
    var msg = "Tradier - DoAuthCallback - No auth code provided. (#2)"
    services.Error(errors.New(msg), msg)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return 
  }
  
  // Request and get an access token.
	data := strings.NewReader("grant_type=authorization_code&code=" + code)
	
	req, err := http.NewRequest("POST", apiBaseUrl + "/oauth/accesstoken", data)
	
	if err != nil {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#1)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return
	}
	
	req.SetBasicAuth(os.Getenv("TRADIER_CONSUMER_KEY"), os.Getenv("TRADIER_CONSUMER_SECRET"))
	req.Header.Add("Accept", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)	
	
  if err != nil {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#2)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return
	}
	
	defer resp.Body.Close()

  // Make sure we got a good status code	
	if resp.StatusCode != http.StatusOK {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#3)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return
	}
	
	// Get the json out of the body.
	jsonBody, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#4)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return
	}
  
  // Put json into an object.
  var tr tokenResponse
  
  err = json.Unmarshal(jsonBody, &tr)
	
	if err != nil {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#5)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return
	}
	
	// Make sure this request was approved.
	if tr.Status != "approved" {
    services.Error(err, "Tradier - DoAuthCallback - Failed to get access token. (#6)")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)
		return  	
	}
	
  // Start the database
  var DB = models.DB{}
  DB.Start()
  defer DB.Connection.Close()  	
  
  // Make sure this is a valid user.
  u, _ := strconv.ParseUint(state, 10, 32)
  user, err := DB.GetUserById(uint(u))
  
  if err != nil {
    var msg = "Tradier - DoAuthCallback - No user found."
    services.Error(err, msg)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return
  }  
  
  // Create new broker entry.
  _, err2 := DB.CreateNewBroker("Tradier", user, tr.Token, tr.RefreshToken, time.Now().Add(time.Duration(tr.ExpiresSec) * time.Second).UTC())

  if err2 != nil {
    var msg = "Tradier - DoAuthCallback - Failed to create broker."
    services.Error(err2, msg)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(genericError)   
    return
  }  

  // Log
  services.Log("Tradier authorization completed for " + user.Email)
 
  // Return success redirect
  http.Redirect(w, r, "/", 302)
  
}

//
// Check to see if we need to refresh the refresh token.
//
func (t * Api) DoRefreshAccessTokenIfNeeded(user models.User) error {

  // Start the database
  var DB = models.DB{}
  DB.Start()
  defer DB.Connection.Close()  
    
  // Get the different tradier brokers.
  brokers, err := DB.GetBrokerTypeAndUserId(user.Id, "Tradier")

  if err != nil {
    services.Error(err, "Tradier - DoRefreshAccessTokenIfNeeded - No brokers found.")
    return err
  }
  
  // Loop through and deal with each tradier broker in the db.
  for i, _ := range brokers {
    
    // Is it time to refresh
    if time.Now().UTC().Add(1 * time.Hour).After(brokers[i].TokenExpirationDate.UTC()) {
      
      err, msg := DoRefreshAccessToken(DB, brokers[i])
      
      if err == nil {
        
        // Update the access token.
        t.muApiKey.Lock()
        t.ApiKey = msg
        t.muApiKey.Unlock()
        
        services.Log("Refreshed Tradier token : " + user.Email)
        
      } else {
        
        services.Error(err, msg + " : " + user.Email)        
      
      }
      
    }    
        
  }
  
  // All done no errors
  return nil
  
}

//
// Get a new access token via the refresh token.
//
func DoRefreshAccessToken(DB models.DB, broker models.Broker) (error, string) {
  
  // Decrypt the refresh token
  decryptRefreshToken, err := helpers.Decrypt(broker.RefreshToken) 
  
  if err != nil {
    return err, "Tradier - DoRefreshAccessToken - FUnable to decrypt message (#1)"
  }    
  
  // Request and get an access token.
	data := strings.NewReader("grant_type=refresh_token&refresh_token=" + decryptRefreshToken)
	
	req, err := http.NewRequest("POST", apiBaseUrl + "/oauth/refreshtoken", data)
	
	if err != nil {
    return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#1)"
	}
	
	req.SetBasicAuth(os.Getenv("TRADIER_CONSUMER_KEY"), os.Getenv("TRADIER_CONSUMER_SECRET"))
	req.Header.Add("Accept", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)	
	
  if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#2)"
	}
	
	defer resp.Body.Close()

  // Make sure we got a good status code	
	if resp.StatusCode != http.StatusOK {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#3)"
	}
	
	// Get the json out of the body.
	jsonBody, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#4)"
	}  
  
  // Put json into an object.
  var tr tokenResponse
  
  err = json.Unmarshal(jsonBody, &tr)
	
	if err != nil {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#5)"
	}
	
	// Make sure this request was approved.
	if tr.Status != "approved" {
		return err, "Tradier - DoRefreshAccessToken - Failed to get access token. (#6)" 	
	}  
  
  // Update the database
  broker.AccessToken = tr.Token
  broker.RefreshToken = tr.RefreshToken
  broker.TokenExpirationDate = time.Now().Add(time.Duration(tr.ExpiresSec) * time.Second).UTC()
  DB.UpdateBroker(broker)

  // All done no errors
  return nil, tr.Token  
  
}

/* End File */