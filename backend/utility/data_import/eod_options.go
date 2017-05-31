//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import(
	"os"
	"fmt" 
  "github.com/tj/go-dropy"
  "github.com/tj/go-dropbox"
  "app.options.cafe/backend/library/services"  
)

//
// Do End of Day Options Import....
//
func DoEodOptionsImport() {
    
  // Get the Dropbox Client (this is where we archive the zip file)
  client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
  
  // Get files we should skip.
  skip, err := GetProccessedFiles(client)
  
  if err != nil {
    services.Error(err, "Could not get skipped files in getProccessedFiles()")
    return
  }  
  
  for key, _ := range skip {
    
    fmt.Println(key)
    
  }
  
  //client.Upload("/demo.txt", strings.NewReader("Hello World"))

}

//
// Get a map of files we have already proccessed.
//
func GetProccessedFiles(client *dropy.Client) (map[string]os.FileInfo, error) {
  
  // Files that we already have imported
  skip := make(map[string]os.FileInfo)  
  
  // Get a list of all the EOD zip files from Dropbox
  files, err := client.List("/data/AllOptions/Daily")
  
  if err != nil {
    return nil, err
  }
  
  // Create a map (hash table) of all the files we already have at Dropbox
  for _, row := range files {
    skip[row.Name()] = row    
  }
  
  return skip, nil
  
}

/* End File */