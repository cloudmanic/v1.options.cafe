package controllers

import (
	"net/http"
)

//
// Return the html tmplate of app.
//
func (t *Controller) HtmlMainTemplate(w http.ResponseWriter, r *http.Request) {

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
      var app_server = "https://app.options.cafe";
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
