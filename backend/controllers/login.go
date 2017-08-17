package controllers

import (
  "net/http"
  "encoding/json" 
  "app.options.cafe/backend/library/realip"
  "app.options.cafe/backend/library/services"     
)

//
// Login to account.
//
func DoLogin(w http.ResponseWriter, r *http.Request) {
    
  // Make sure this is a post request.
	if r.Method == http.MethodGet {
    HtmlMainTemplate(w, r)
    return
	} 
	
  // Make sure this is a post request.
	if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
	}
	
	// Set response
	w.Header().Set("Content-Type", "application/json")
  
  // Decode json passed in
  decoder := json.NewDecoder(r.Body)
  
  type LoginPost struct {  
    Email string
    Password string
  }
  
  var post LoginPost 
  
  err := decoder.Decode(&post)
  
  if err != nil {
    services.Error(err, "DoLogin - Failed to decode JSON posted in")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))   
    return 
  }
  
  defer r.Body.Close()
  
  // Start the database
  DB.Start()
  defer DB.Connection.Close()

  // Validate user.
  if err := DB.ValidateUserLogin(post.Email, post.Password); err != nil {
    
    // Respond with error
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"" + err.Error() + "\"}"))     
    
    return 
  }  

  // Login user in by email and password
  user, err := DB.LoginUserByEmailPass(post.Email, post.Password, r.UserAgent(), realip.RealIP(r))

  if err != nil {
    services.Error(err, "DoLogin - Unable to log user in. (CreateUser)")
    
    // Respond with error
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Sorry, we could not find an account with that email / password combination.\"}"))     
    
    return     
  }  
  
  // Here we check to see if we have any brokers. If there are no brokers the user needs to select at least one to do anything.
  var brokerCount = len(user.Brokers)
  
  type Response struct {
    Status uint `json:"status"`
    AccessToken string `json:"access_token"`
    BrokerCount int `json:"broker_count"`
  }
  
  resObj := &Response{ 
    Status: 1,
    AccessToken: user.Session.AccessToken,
    BrokerCount: brokerCount,
  }

  resJson, err := json.Marshal(resObj)
  
  if err != nil {
    services.Error(err, "DoLogin - Unable to log user in. (json.Marshal)") 
    
    // Respond with error
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))     
    
    return     
  } 

  // Return success json.
  w.Write(resJson) 
}

/* End File */