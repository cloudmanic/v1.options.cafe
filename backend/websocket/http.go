package websocket

import (
  "os"
  "fmt"
  "log"
  "time"
  "net/http"
  "golang.org/x/crypto/acme/autocert"
)

var (
  WsWriteChannel chan SendStruct
  WsWriteQuoteChannel chan SendStruct
)

//
// Start the webserver
//
func Start() {
  
  // Listen for data from our broker feeds.
  go DoWebsocketDataFeeds()
  
  // Register some handlers:
  mux := http.NewServeMux()
  
  // Http Routes
  mux.HandleFunc("/", HtmlMainTemplate)    
    
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
func DoWebsocketDataFeeds() {
  
  for {
  
    select {
    
      case send := <-WsWriteChannel:
        fmt.Println(send.Message)
        
      case send := <-WsWriteQuoteChannel:
        fmt.Println(send.Message)
      
    }
  
  }   
  
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