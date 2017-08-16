package controllers

import (
  "fmt"
  "net/http"
  "encoding/json" 
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
    w.Header().Set("Content-Type", "application/json")
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
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"" + err.Error() + "\"}"))     
    
    return 
  }  

  // Login user in by email and password
  user, err := DB.LoginUserByEmailPass(post.Email, post.Password, r.UserAgent(), r.RemoteAddr)

  if err != nil {
    services.Error(err, "DoLogin - Unable to log user in. (CreateUser)")
    
    // Respond with error
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))     
    
    return     
  }  

  // Return success json.
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("{\"status\":1, \"access_token\":\"" + user.Session.AccessToken + "\"}"))  
}

/* End File */