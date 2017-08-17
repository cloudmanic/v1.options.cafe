package controllers

import (
  "net/http"
  "encoding/json"
  "app.options.cafe/backend/library/realip" 
  "app.options.cafe/backend/library/services"     
)

//
// Register a new account.
//
func DoRegister(w http.ResponseWriter, r *http.Request) {
    
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
  
  type RegisterPost struct {
    First string
    Last string    
    Email string
    Password string
  }
  
  var post RegisterPost 
  
  err := decoder.Decode(&post)
  
  if err != nil {
    services.Error(err, "DoRegisterPost - Failed to decode JSON posted in")
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong while registering your account. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))   
    return 
  }
  
  defer r.Body.Close()
  
  // Start the database
  DB.Start()
  defer DB.Connection.Close()

  // Validate user.
  if err := DB.ValidateCreateUser(post.First, post.Last, post.Email, post.Password); err != nil {
    
    // Respond with error
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"" + err.Error() + "\"}"))     
    
    return 
  }  

  // Install new user.
  user, err := DB.CreateUser(post.First, post.Last, post.Email, post.Password, r.UserAgent(), realip.RealIP(r))

  if err != nil {
    services.Error(err, "DoRegisterPost - Unable to register new user. (CreateUser)")
    
    // Respond with error
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("{\"status\":0, \"error\":\"Something went wrong while registering your account. Please try again or contact help@options.cafe. Sorry for the trouble.\"}"))     
    
    return     
  }

  // Return success json.
  w.Write([]byte("{\"status\":1, \"access_token\":\"" + user.Session.AccessToken + "\"}"))  
}

/* End File */