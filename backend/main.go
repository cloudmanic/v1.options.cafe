package main

import (
  "fmt"
  "log"
  "flag"
  "sync"
  "time"
  "runtime"
  "net/http"
  "./library/https"
  "./brokers/tradier"
  "github.com/gorilla/websocket"
  "golang.org/x/crypto/acme/autocert"
)

var ( 
  mu sync.Mutex
  localDevMode = flag.Bool("local", false, "localhost development")

  
  ws Websockets
  brokerFeeds = make(map[string]*BrokerFeed)
  
  // Channels
  websocketSendChannel = make(chan string)
  websocketSendQuoteChannel = make(chan string)
)
    
//
// Main....
//
func main() {  
  
  // Parse flags
  flag.Parse()
  
  // Lets get started
  println("App Started.....")
 
  // Setup CPU stuff.
  runtime.GOMAXPROCS(runtime.NumCPU())  
    
  // Setup websocket connections.
  ws.connections = make(map[*websocket.Conn]*WebsocketConnection)
  ws.quotesConnections = make(map[*websocket.Conn]*WebsocketConnection)
  
  // Get started websocket sending
  go ws.DoWebsocketSending()
  go ws.DoWebsocketQuoteSending()
    
  // Register some handlers:
  mux := http.NewServeMux()
  
  mux.HandleFunc("/", HtmlMainTemplate)    
    
  // Setup websocket
	mux.HandleFunc("/ws/core", ws.DoWebsocketConnection)
	mux.HandleFunc("/ws/quotes", ws.DoQuoteWebsocketConnection)

  // Are we in testing mode?
  if *localDevMode {
    
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
    https.StartSecureServer(mux, m.GetCertificate)
    
  }
	  
}

//
// Start a broker feed. One per lisc.
//
func StartBrokerFeed(licenseKey string, brokerApiToken string) {
  
  // Lock the memory
	mu.Lock()
	defer mu.Unlock() 

  // Make sure we do not already have this licenseKey going.
  if _, ok := brokerFeeds[licenseKey]; ok {
    return
  }  

  // Log we are starting this.
  fmt.Println("Starting The Broker Feed - " + licenseKey)
  
  // Set the broker we are going to use.
  var broker = tradier.Api{ ApiKey: brokerApiToken }
  var fetch = Fetch{ broker: broker }
  
  // Set Broker hash and lets get going.
  brokerFeeds[licenseKey] = &BrokerFeed{ 
                              fetch: fetch, 
                              licenseKey: licenseKey, 
                              brokerApiToken: brokerApiToken,
                            }

  // Start the broker feed.
  brokerFeeds[licenseKey].Start(websocketSendQuoteChannel)
}

//
// Set the active broker account id.
//
func SetBrokerAccountId(licenseKey string, brokerAccountId string) {
  
  // Lock the memory
	mu.Lock()
	defer mu.Unlock() 

  // Make sure we do not already have this licenseKey going.
  if _, ok := brokerFeeds[licenseKey]; ! ok {
    return
  }  
  
  // Set the account id.
  brokerFeeds[licenseKey].brokerAccountId = brokerAccountId
  brokerFeeds[licenseKey].fetch.broker.SetDefaultAccountId(brokerAccountId)
  
	// Log this action.
	fmt.Println("New Default Account Id From Frontend - " + brokerAccountId)  
  
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
  
  	<link rel="shortcut icon" type="image/x-icon" href="assets/css/images/favicon.ico?v=1" />
  	<link rel="icon" type="image/png" href="assets/css/images/favicon-32x32.png?v=1" sizes="32x32" />
  	<link rel="icon" type="image/png" href="assets/css/images/favicon-16x16.png?v=1" sizes="16x16" />
  	
  	<link rel="stylesheet" href="assets/vendor/bootstrap-3.3.7-dist/css/bootstrap.min.css?v=1" type="text/css" media="all" />
  	<link rel="stylesheet" href="assets/css/style.css?v=1" />  
  	
    <script type="text/javascript">
      
      // Set websocket url
      var ws_server = "wss://app.options.cafe";
      //var ws_server = "ws://127.0.0.1:7652";
      
      // Set tradier api key
      var tradier_api_key = localStorage.getItem('tradier_api_key');
    </script>  		
  </head>
<body>
  <div class="wrapper">
    <oc-root>Loading...</oc-root>
  </div>

  <script src="assets/vendor/jquery-1.12.4.min.js?v=1"></script>
  <script src="assets/vendor/bootstrap-3.3.7-dist/js/bootstrap.min.js?v=1"></script>
  <script src="assets/js/functions.js?v=1"></script>
  <script type="text/javascript" src="inline.bundle.js?v=1"></script>
  <script type="text/javascript" src="vendor.bundle.js?v=1"></script>
  <script type="text/javascript" src="main.bundle.js?v=1"></script>        
</body>
</html>  
  `))
  
}

/* End File */
