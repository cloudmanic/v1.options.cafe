//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import(
	"os"
  "io"
  "net/http"
  "github.com/tj/go-dropy"
  "github.com/tj/go-dropbox"
  "github.com/jlaffaye/ftp"
  "app.options.cafe/backend/library/services"  
)

//
// Do End of Day Options Import....
//
func DoEodOptionsImport() {

  // Log
  services.Log("Starting DoEodOptionsImport().")

  // Get the Dropbox Client (this is where we archive the zip file)
  client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
  
  // Get files we should skip.
  dbFiles, err := GetProccessedFiles(client)
  
  if err != nil {
    services.Error(err, "Could not get dbFiles files in GetProccessedFiles()")
    return
  }  

  ftpFiles, err := GetOptionsDailyData()

  if err != nil {
    services.Error(err, "Could not get ftp files in GetOptionsDailyData()")
    return
  }    

  for _, row := range ftpFiles {

    if _, ok := dbFiles[row.Name]; ok {
      continue;
    }

    // Download file.
    filePath, err := DownloadOptionsDailyDataByName(row.Name)

    if err != nil {
      services.Error(err, "Could not download ftp file in DownloadOptionsDailyDataByName() - " + row.Name)
      continue;
    }

    // Open file we upload.
    file, err := os.Open(filePath)

    if err != nil {
      services.Error(err, "Could not open ftp file os.Open() - " + row.Name)
      continue;      
    }

    // Upload the file to Dropbox
    client.Upload("/data/AllOptions/Daily/" + row.Name, file)

    if err != nil {
      services.Error(err, "Could not upload to Dropbox - " + row.Name)
      continue;      
    }

    // Log file
    services.Log("Finished uploading " + row.Name + " to Dropbox.")

    // Delete file we uploaded to Dropbox
    err = os.Remove(filePath)

    if err != nil {
      services.Error(err, "Could not delete file - " + filePath)
      continue;      
    }

  }

  // Send health check notice.
  if len(os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL")) > 0 {

    resp, err := http.Get(os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL"))
    
    if err != nil {
      services.Error(err, "Could send health check - " + os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL"))
    }
    
    defer resp.Body.Close()
    
  }

  // Log 
  services.Log("Done DoEodOptionsImport().")  

}

//
// Download file from Delta Neutral
//
func DownloadOptionsDailyDataByName(name string) (string, error) {

  client, err := ftp.Dial("eodftp.deltaneutral.com:21")

  defer client.Quit()

  if err != nil {
    return "", err
  }

  if err := client.Login(os.Getenv("DELTA_NEUTRAL_USERNAME"), os.Getenv("DELTA_NEUTRAL_PASSWORD")); err != nil {
    return "", err
  }

  reader, err := client.Retr("/dbupdate/" + name)

  if err != nil {
    return "", err
  }

  // Save file locally
  writer, err := os.Create("/tmp/" + name)

  if err != nil {
    return "", err
  }

  if _, err = io.Copy(writer, reader); err != nil {
    return "", err
  }

  return "/tmp/" + name, nil

}

//
// Get FTP options data from provider.
//
func GetOptionsDailyData() ([]*ftp.Entry, error) {

  client, err := ftp.Dial("eodftp.deltaneutral.com:21")

  defer client.Quit()

  if err != nil {
    return nil, err
  }

  if err := client.Login(os.Getenv("DELTA_NEUTRAL_USERNAME"), os.Getenv("DELTA_NEUTRAL_PASSWORD")); err != nil {
    return nil, err
  }

  entries, _ := client.List("/dbupdate/options_*.zip")

  return entries, nil
}

//
// Get a map of files we have already processed.
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