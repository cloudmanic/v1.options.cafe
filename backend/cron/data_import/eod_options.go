//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import(
	"os"
  "io"
  "fmt"
  "time"
  "errors"
  "strings"
  "net/http"
  "encoding/csv"
  "github.com/tj/go-dropy"
  "github.com/tj/go-dropbox"
  "github.com/jlaffaye/ftp"
  "github.com/araddon/dateparse"
  "app.options.cafe/backend/library/services"  
)

// Job Struct
type Job struct {
  Key string
  Length int
  Index int
  Path string
}

//
// Do End of Day Options Import. We run this every day from "Cron". It will 
// connect to the FTP site at DeltaNeutral (our EOD options data provider) and 
// download any days data that we do not already have. Once we download the file we
// upload it to Dropbox for achieving. Then we break up the DeltaNeutral export
// into one file per date per asset symbol. Lastly we upload this data to AWS. 
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

  // Get daily options from vendor
  ftpFiles, err := GetOptionsDailyData()

  if err != nil {
    services.Error(err, "Could not get ftp files in GetOptionsDailyData()")
    return
  }    

  // Loop through the FTP files.
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

    // Import file into different symbols
    err = SymbolImport(filePath)

    if err != nil {
      services.Error(err, "Could not import DatabaseImport() - " + row.Name)
      continue;      
    }

    // // Delete file we uploaded to Dropbox
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
// Convert all Delta Neutral imports into a per symbol / date CSV file. This is a mass import 
// from Dropbox where we archive Delta Neutral EOD options data. We go through each 
// EOD file and convert to a per symbol per date archive. We then store this archive at 
// AWS S3. This is used as a one off thing when data at AWS gets messed up or is missing or something.
// In a perfect world this is only run once ever.
//
func ProccessAllDeltaNeutralData() error {

  // Log
  services.Log("Starting ProccessAllDeltaNeutralData().")

  // Job queue stuff. (assume we do not have more than 2000 days worth of data)
  jobs := make(chan Job, 2000)
  results := make(chan string, 2000)

  // Get the Dropbox Client (this is where we archive the zip file)
  client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
  
  // Get files we should skip.
  dbFiles, err := GetProccessedFiles(client)  

  if err != nil {
    return err
  }

  // Start workers
  for w := 1; w <= 50; w++ {
    go ProccessDeltaNeutralDataWorker(w, jobs, results)
  }  

  var i = 1;
  var total = len(dbFiles)

  // Loop through files at Dropbox
  for key, _ := range dbFiles {

    // Send job to worker
    jobs <- Job{ Key: key, Index: i, Length: total, Path: "/data/AllOptions/Daily/" + key }

    // Update count
    i++
  }   

  // Close down jobs
  close(jobs)

  // Collect the results of the workers
  for a := 0; a < total; a++ {
    <-results
  }

  // Close down results
  close(results)

  // Log 
  services.Log("Done ProccessAllDeltaNeutralData().")   

  // Return happy.
  return nil
}

//
// A worker to process jobs of DeltaNeutral imports
//
func ProccessDeltaNeutralDataWorker(id int, jobs <-chan Job, results chan<- string) {

  // Loop through the different jobs until we have no more jobs.
  for job := range jobs {

    fmt.Printf("Downloading: %s (%d/%d) (worker #%d) \n", job.Path, job.Index, job.Length, id)

    client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))

    file, err := client.Download(job.Path);    

    if err != nil {
      services.Error(err, "Could not download from Dropbox. ProccessDeltaNeutralDataWorker()")
      return
    }

    // Create Zip file to store
    outFile, err := os.Create("/Volumes/Drive04/options/" + job.Key);

    if err != nil {
      services.Error(err, "Could not create file. ProccessDeltaNeutralDataWorker()")      
      return
    }  

    // Copy DB file downloaded.
    _, err = io.Copy(outFile, file)

    if err != nil {
      services.Error(err, "Could not copy from Dropbox. ProccessDeltaNeutralDataWorker()")      
      return
    }

    // Write file.
    outFile.Close()

    // Import symbol
    SymbolImport("/Volumes/Drive04/options/" + job.Key)

    // Delete file.
    os.Remove("/Volumes/Drive04/options/" + job.Key) 

    // Send result back so we know it is done.
    results <- job.Key
  }

  // Log and return 
  fmt.Printf("Worker %d closed.", id)

  // Return happy
  return
}

//
// Break the data up per symbol
//
func SymbolImport(filePath string) error {

  // Log
  services.Log("Start SymbolImport - " + filePath) 

  // Unzip CSV files.
  files, err := Unzip(filePath, "/tmp/output/")

  if err != nil {
    return err
  }

  // Loop through the files paths and find the file we are looking for.
  var file = ""

  for _, row := range files {

    // We are only interested in one file.
    i := strings.Index(row, "ptions_")    

    if i > -1 {
      file = row
    }

  }

  // Make sure we have a file
  if len(file) <= 0 {
    return errors.New("Did not find /tmp/output/options_XXXXXXXX.csv")
  }

  // Open CSV file
  f, err := os.Open(file)
  
  if err != nil {
    return err
  }
  
  defer f.Close()

  // Read File into a Variable - https://golangcode.com/how-to-read-a-csv-file-into-a-struct
  lines, err := csv.NewReader(f).ReadAll()
  
  if err != nil {
    return err
  }

  // Log import
  services.Log("Importing option EOD quotes for - " + string(lines[0][7]))

  // Figure out quote date
  date, err := dateparse.ParseAny(string(lines[0][7]))

  if err != nil {
    return err
  }

  // Hash map to keep each symbol array
  symbolMap := make(map[string][][]string)

  // Loop through lines & turn into object
  for _, line := range lines {
    symbolMap[line[0]] = append(symbolMap[line[0]], line)
  } 

  // Loop through and store to file.
  for key, _ := range symbolMap {  

    // Store Symbol to a file based on date and symbol
    zipFilePath, err := StoreOneDaySymbol(key, date, symbolMap[key])
    
    if err != nil {
      return err
    }

    // Send the file to AWS for storage
    err = AWSUpload(zipFilePath, key)

    if err != nil {
      return err
    }

  }

  // Log
  services.Log("Done SymbolImport - " + filePath)   

  // Return happy.
  return nil
}

//
// Store one day's worth of Symbols 
//
func StoreOneDaySymbol(symbol string, date time.Time, data [][]string) (string, error) {

  var fileName = date.Format("2006-01-02") + ".csv"
  var dirBase = os.Getenv("CACHE_DIR") + "/options-eod/"
  var dirPath = dirBase + symbol + "/"
  var csvFile = dirPath + fileName  
  var zipFile = csvFile + ".zip"

  // Create directory - Base
  if _, err := os.Stat(dirBase); os.IsNotExist(err) {
    err = os.Mkdir(dirBase, 0755)

    if err != nil {
      return "", err
    }     
  }  

  // Create directory - Path
  if _, err := os.Stat(dirPath); os.IsNotExist(err) {
    err = os.Mkdir(dirPath, 0755)

    if err != nil {
      return "", err
    }     
  } 

  // Create CSV file
  csvFilePtr, err := os.Create(dirPath + fileName);

  if err != nil {
    return "", err
  }  

  // Write to a new CSV file just for this symbol
  writer := csv.NewWriter(csvFilePtr)
    
  // Loop through writing each line to the file.
  for _, row := range data {

    // Make sure the row is not blank
    if err := writer.Write(row); err != nil {
      return "", err
    }

  }  

  // Write the file.
  writer.Flush()
  csvFilePtr.Close()

  // Zip the file up.
  err = ZipFiles(zipFile, []string{csvFile})  

  if err != nil {
    return "", err
  }

  // Delete unziped file
  err = os.Remove(csvFile)

  if err != nil {
    return "", err      
  }  

  // Return happy
  return zipFile, err
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

  // Return file path
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


/*
    // Delta Neutral Options CSV Format
    [0] => underlying
    [1] => underlying_last
    [2] =>  exchange
    [3] => optionroot
    [4] => optionext
    [5] => type
    [6] => expiration
    [7] => quotedate
    [8] => strike
    [9] => last
    [10] => bid
    [11] => ask
    [12] => volume
    [13] => openinterest
    [14] => impliedvol
    [15] => delta
    [16] => gamma
    [17] => theta
    [18] => vega
    [19] => optionalias
*/

/* End File */