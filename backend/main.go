package main

import (
  "os"
  "fmt"
  "log"
  "sync"
  "time"
  "runtime"
  "net/http"
  "github.com/stvp/rollbar"
  "github.com/jinzhu/gorm"
  "github.com/joho/godotenv"
  "github.com/gorilla/websocket"
  "golang.org/x/crypto/acme/autocert"
  _ "github.com/go-sql-driver/mysql"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/library/https"
  "app.options.cafe/backend/brokers/tradier"  
)

var ( 
  mu sync.Mutex
  db *gorm.DB
  
  // Websocket connections.
  ws = Websockets{      
    connections: make(map[*websocket.Conn]*WebsocketConnection),
    quotesConnections: make(map[*websocket.Conn]*WebsocketConnection),
    controller: WsController{},
  }
  
  // User connections
  userConnections = make(map[uint]*UsersConnection)
)

type UsersConnection struct {
  UserId uint
  BrokerConnections map[uint]*BrokerFeed
  WsWriteChannel chan string
  WsWriteQuoteChannel chan string
}
        
//
// Main....
//
func main() {  
   
  // Setup CPU stuff.
  runtime.GOMAXPROCS(runtime.NumCPU())  
         
  // Load .env file 
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }        
    
  // Setup Rollbar
  rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
  rollbar.Environment = os.Getenv("ROLLBAR_ENV")    
    
  // Lets get started
  fmt.Println("App Started: " + os.Getenv("APP_ENV"))    
  
  // Message that the app has started
  rollbar.Message("info", "App started.")
    
  // Connect to database and run Migrations.
  db = DbConnect()

  // Start our broker connections.
  StartUserConnections()
    
  // Register some handlers:
  mux := http.NewServeMux()
  
  // Http Routes
  mux.HandleFunc("/", HtmlMainTemplate)    
    
  // Setup websocket
	mux.HandleFunc("/ws/core", ws.DoWebsocketConnection)
	mux.HandleFunc("/ws/quotes", ws.DoQuoteWebsocketConnection)

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
    https.StartSecureServer(mux, m.GetCertificate)
    
  }
	  
}

//
// Start up our user connections.
//
func StartUserConnections() {
  
  var user models.User
  users := user.GetAllUsers(db)

  for i, _ := range users {
    go StartUserConnection(users[i])     
  }  
  
}

//
// Start one user connection.
//
func StartUserConnection(user models.User) {
  
  fmt.Println("Starting User Connection : " + user.Email)
  
  // Lock the memory
	mu.Lock()
	defer mu.Unlock() 

  // Make sure we do not already have this licenseKey going.
  if _, ok := userConnections[user.Id]; ok {
    rollbar.Message("info", "User Connection Is Already Going : " + user.Email)
    fmt.Println("User Connection Is Already Going : " + user.Email)
    return
  }
  
  // Set the user connection.
  userConnections[user.Id] = &UsersConnection{
    UserId: user.Id,
    BrokerConnections: make(map[uint]*BrokerFeed),
    WsWriteChannel: make(chan string),
    WsWriteQuoteChannel: make(chan string),
  }
  
  // Start the websocket write dispatcher for this user.
  go ws.DoWsDispatch(userConnections[user.Id])
  
  // Loop through the different brokers for this user
  for _, row := range user.Brokers {
    
    // Need an access token to continue
    if len(row.AccessToken) <= 0 {
      rollbar.Message("info", "User Connection (Brokers) No Access Token Found : " + user.Email)
      fmt.Println("User Connection (Brokers) No Access Token Found : " + user.Email)
      continue
    }
    
    // Set the broker we are going to use.
    var broker = tradier.Api{ ApiKey: row.AccessToken }
    var fetch = Fetch{ broker: broker, user: userConnections[user.Id] }
    
    // Set Broker hash and lets get going.
    userConnections[user.Id].BrokerConnections[row.Id] = &BrokerFeed{ fetch: fetch, userId: user.Id }
  
    // Start the broker feed.
    userConnections[user.Id].BrokerConnections[row.Id].Start()     
    
  }
    
}

//
// Connect to the db and run migrations.
//
func DbConnect() (*gorm.DB) {
  
  var err error
    
  // Connect to Mysql
  conn, err := gorm.Open("mysql", os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + "/" + os.Getenv("DB_DATABASE") + "?charset=utf8&parseTime=True&loc=Local")
  
  if err != nil {
    rollbar.Error(rollbar.ERR, err)
    panic("failed to connect database")
  }

  // Enable
  //conn.LogMode(true)
  //conn.SetLogger(log.New(os.Stdout, "\r\n", 0))

  // Migrate the schemas (one per table).
  conn.AutoMigrate(&models.User{})
  conn.AutoMigrate(&models.Broker{})
  conn.AutoMigrate(&models.Order{})
  conn.AutoMigrate(&models.OrderLeg{})
  conn.AutoMigrate(&models.Position{}) 
  conn.AutoMigrate(&models.TradeGroup{}) 
  
  return conn   
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

/* End File */
