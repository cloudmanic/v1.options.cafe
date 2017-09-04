//
// Date: 8/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
  "os"
  "bytes"
  "errors"
  "net/http"
  "io/ioutil"
)

//
// Slack notify
//
func SlackNotify(channel string, msg string) (string, error) {
  
  if len(os.Getenv("SLACK_HOOK")) > 0 {
  
    var jsonStr = []byte(`{"channel": "` + channel + `", "text": "` + msg + `"}`)
  
    // Creatre POST request  
    req, err := http.NewRequest("POST", os.Getenv("SLACK_HOOK"), bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
  
    // Send request.
    client := &http.Client{}
    resp, err := client.Do(req)
    
    if err != nil {
      Error(err, "SlackNotify - Unable to send slack notice : " + msg + ".")
      return "", err
    }
    
    if resp.StatusCode != http.StatusOK {
      Error(err, "SlackNotify (no 200) - Unable to send slack notice : " + msg + ".")
      return "", err          
    }  
  
    // Get the body.
    body, _ := ioutil.ReadAll(resp.Body)    
  
    resp.Body.Close()
    
    // Return happy.
    return string(body), err 
  
  } 
  
  // Nothing happened.
  return "", errors.New("SLACK_HOOK is not set.")
  
}

/* End File */