package websocket

import (
  "os"
  "log"
  "time"
  "net/http"
  "encoding/json"
  "golang.org/x/crypto/acme/autocert"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/library/services"
)

var (
  DB *models.DB
  WsReadChan chan SendStruct
  WsWriteChan chan SendStruct
  WsWriteQuoteChan chan SendStruct
)

//
// Start the webserver
//
func Start() {
  
  // Listen for data from our broker feeds.
  go DoWsDispatch()
  
  // Register some handlers:
  mux := http.NewServeMux()
  
  // Http Routes
  mux.HandleFunc("/", HtmlMainTemplate)    
  mux.HandleFunc("/register", DoRegisterPost) 
    
  // Setup websocket
	mux.HandleFunc("/ws/core", DoWebsocketConnection)
	mux.HandleFunc("/ws/quotes", DoQuoteWebsocketConnection)

  // Are we in testing mode?
  if os.Getenv("APP_ENV") == "local" {
    
		s := &http.Server{
			Addr: ":7652",
			Handler: mux,
      ReadTimeout:  2 * time.Second,
      WriteTimeout: 2 * time.Second,			
		}
		
		log.Fatal(s.ListenAndServe())    
     
  } else {

    // Secure it with a TLS certificate using Let's  Encrypt:
    m := autocert.Manager{
  	  Prompt: autocert.AcceptTOS,
      Cache: autocert.DirCache("/etc/letsencrypt/"),
      Email: "support@options.cafe",
      HostPolicy: autocert.HostWhitelist("app.options.cafe"),
    }
  
    // Start a secure server:
    StartSecureServer(mux, m.GetCertificate)
    
  }
  
}

//
// Listen for data from our broker feeds.
// Take the data and then pass it up the websockets.
//
func DoWsDispatch() {
  
  for {
    
    select {

      // Core channel
      case send := <-WsWriteChan:
      
        for i, _ := range connections {
          
          // We only care about the user we passed in.
          if connections[i].userId == send.UserId {
            
            select {
              
              case connections[i].writeChan <-send.Message:
 	
              default:
                services.MajorLog("Channel full. Discarding value (Core channel)")   
                          
            }
            
          }
          
        }
      
      // Quotes channel
      case send := <-WsWriteQuoteChan:
      
        for i, _ := range quotesConnections {
          
          // We only care about the user we passed in.
          if quotesConnections[i].userId == send.UserId {
            
             select {
              
              case quotesConnections[i].writeChan <-send.Message:
 	
              default:
                services.MajorLog("Channel full. Discarding value (Quotes channel)")   
                          
            }           
            
          }
          
        }
         
    }
      
  }  
    
}

//
// Register a new account.
//
func DoRegisterPost(w http.ResponseWriter, r *http.Request) {
    
  // Make sure this is a post request.
	if r.Method != "POST" {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
	} 
  
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
    return 
  }
  
  defer r.Body.Close()
  
  // Start the database
  DB.Start()
  defer DB.Connection.Close()
  
  // First see if we already have this user.
  user, err := DB.GetUserByEmail(post.Email)
  
  if err == nil {
    services.Log("DoRegisterPost - Attempt to register user already in database : " + user.Email)
    
    // Respond with error
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{\"status\":0, \"error\":\"Looks like you already have an account, please login?\"}"))    
      
    return
  }

  // Install new user.
  user, err = DB.CreateUser(post.First, post.Last, post.Email, post.Password)

  if err != nil {
    services.Error(err, "DoRegisterPost - Unable to register new user. (CreateUser)")
    return     
  }

  // Return success json.
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("{\"status\":1, \"access_token\":\"" + user.AccessToken + "\"}"))  
}


//
// Return the html tmplate of app.
//
func HtmlMainTemplate(w http.ResponseWriter, r *http.Request) {

  w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
	<head>
  	<base href="https://cdn.options.cafe/app/" />
  	<meta charset="utf-8" />
  
  	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no" />
  
  	<title>Options Cafe</title>
  
  	<link rel="shortcut icon" type="image/x-icon" href="assets/css/images/favicon.ico?v=9" />
  	<link rel="icon" type="image/png" href="assets/css/images/favicon-32x32.png?v=9" sizes="32x32" />
  	<link rel="icon" type="image/png" href="assets/css/images/favicon-16x16.png?v=9" sizes="16x16" />
  	
  	<link rel="stylesheet" href="assets/vendor/bootstrap-3.3.7-dist/css/bootstrap.min.css" type="text/css" media="all" />
  	<link rel="stylesheet" href="assets/css/style.css?v=9" />  
  	
    <script type="text/javascript">
      var ws_server = "wss://app.options.cafe";
    </script>  		
  </head>
<body>
  <div class="wrapper">
    <oc-root>Loading...</oc-root>
  </div>

  <script src="assets/vendor/jquery-1.12.4.min.js"></script>
  <script src="assets/vendor/bootstrap-3.3.7-dist/js/bootstrap.min.js"></script>
  <script src="assets/bower/clientjs/dist/client.min.js"></script>
  <script src="assets/js/functions.js?v=9"></script>
  <script type="text/javascript" src="inline.bundle.js?v=9"></script>
  <script type="text/javascript" src="vendor.bundle.js?v=9"></script>
  <script type="text/javascript" src="main.bundle.js?v=9"></script>        
</body>
</html>  
  `))
  
}