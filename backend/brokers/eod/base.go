//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"strings"
	"time"

	env "github.com/jpfuentes2/go-env"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/cache"
	"app.options.cafe/library/files"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/library/store/object"
	"app.options.cafe/models"
)

const workerCount int = 100
const cacheDirBase = "broker-eod"

var cacheDir string

// Api struct
type Api struct {
	DB  models.Datastore
	Day time.Time // This is the day we pull EOD data for
}

// Job struct
type Job struct {
	Path  string
	Index int
}

// SymbolStore struct
type SymbolStore struct {
	Last    float64                  `json:"last"`
	Options []types.OptionsChainItem `json:"options"`
}

//
// Start up the controller.
//
func init() {
	// Helpful for testing
	env.ReadEnv(build.Default.GOPATH + "/src/app.options.cafe/.env")

	// Set cache dir
	cacheDir = os.Getenv("CACHE_DIR") + "/" + cacheDirBase + "/"

	// Make sure our cache dir is setup.
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}
}

//
// GetOptionsBySymbol - Get options by Symbol. We do not return a chain. We more or less
// return the data in our CSV files as Go structs
//
func (t *Api) GetOptionsBySymbol(symbol string) ([]types.OptionsChainItem, float64, error) {
	// Symbol to upper
	symb := strings.ToUpper(symbol)

	// See if we have this symbol chain cached already. If so return from cache
	o, u, e := getSymbolEodJSONFromCache(symb, t.Day)

	if e == nil {
		return o, u, e
	}

	// Get symbole from S3 store.
	zipFile, err := downloadEodSymbol(symbol, t.Day)

	if err != nil {
		services.Critical(errors.New(err.Error() + "Could not download file - " + zipFile))
	}

	// Convert CSV to Struct
	options, underlyingLast, err := unzipSymbolCSV(symb, zipFile)

	if err != nil {
		return options, underlyingLast, err
	}

	// Store JSON for local cache
	storeSymbolEodJSON(symb, t.Day, underlyingLast, options)

	// Return happy
	return options, underlyingLast, nil
}

//
// Return a list of trade dates by symbol
//
func GetTradeDatesBySymbols(symbol string) ([]time.Time, error) {

	dates := []time.Time{}
	dirPath := "options-eod/" + strings.ToUpper(symbol) + "/"
	cacheKey := "oc-brokers-eod-trade-dates-" + strings.ToUpper(symbol)

	// See if we have this result in the cache.
	cacheddates := []time.Time{}
	found, _ := cache.Get(cacheKey, &cacheddates)

	// Return happy JSON
	if found {
		return cacheddates, nil
	}

	// Make call to S3 store and get the dates for different EODs
	files, err := object.ListObjects(dirPath)

	if err != nil {
		return dates, err
	}

	// Convert strings to dates
	for _, row := range files {
		fileName := strings.Replace(row.Key, dirPath, "", -1)
		fileName = strings.Replace(fileName, ".csv.zip", "", -1)
		dates = append(dates, helpers.ParseDateNoError(fileName).UTC())
	}

	// Store dates in cache - 3 hours
	cache.SetExpire(cacheKey, (time.Minute * 60 * 3), dates)

	// Return happy
	return dates, nil
}

//
// Return a list of keys to download by symbol.
//
func GetTradeDateKeysBySymbol(symbol string) ([]string, error) {

	keys := []string{}
	dirPath := "options-eod/" + strings.ToUpper(symbol) + "/"

	// Get dates from S3 store
	dates, err := GetTradeDatesBySymbols("spy")

	if err != nil {
		return keys, err
	}

	// Loop through and create keys
	for _, row := range dates {
		keys = append(keys, dirPath+row.Format("2006-01-02")+".csv.zip")
	}

	// Return happy
	return keys, nil
}

// ---------------- Private Helper Functions -------------- //

//
// unzipSymbolCSV - Convert CSV lines to objects
//
func unzipSymbolCSV(symb string, zipFile string) ([]types.OptionsChainItem, float64, error) {
	options := []types.OptionsChainItem{}
	underlyingLast := 0.00

	// Make sure we have this file
	if _, err := os.Stat(zipFile); os.IsNotExist(err) {
		return options, underlyingLast, err
	}

	// Unzip option chain
	f, err := files.Unzip(zipFile, "/tmp/"+symb)

	if err != nil {
		return options, underlyingLast, err
	}

	// Open CSV file
	csvFile, err := os.Open(f[0])

	if err != nil {
		return options, underlyingLast, err
	}

	defer csvFile.Close()

	// Read File into a Variable - https://golangcode.com/how-to-read-a-csv-file-into-a-struct
	lines, err := csv.NewReader(csvFile).ReadAll()

	if err != nil {
		return options, underlyingLast, err
	}

	// Set underlyingLast
	underlyingLast = helpers.StringToFloat64(lines[0][1])

	// Loop through the different lines of the CSV and Store in chain
	for _, row := range lines {

		// Get Option parts
		parts, err := helpers.OptionParse(row[3])

		if err != nil {
			return options, underlyingLast, err
		}

		// Build Item
		op := types.OptionsChainItem{
			Underlying:     strings.ToUpper(symb),
			Symbol:         row[3],
			OptionType:     parts.Type,
			Description:    parts.Name,
			Strike:         parts.Strike,
			ExpirationDate: types.Date{parts.Expire},
			Last:           helpers.StringToFloat64(row[9]),
			Volume:         helpers.StringToInt(row[12]),
			Bid:            helpers.StringToFloat64(row[10]),
			Ask:            helpers.StringToFloat64(row[11]),
			OpenInterest:   helpers.StringToInt(row[13]),
			ImpliedVol:     helpers.StringToFloat64(row[11]),
			Delta:          helpers.StringToFloat64(row[15]),
			Gamma:          helpers.StringToFloat64(row[16]),
			Theta:          helpers.StringToFloat64(row[17]),
			Vega:           helpers.StringToFloat64(row[18]),
		}

		// Append Item
		options = append(options, op)
	}

	// Delete csv file
	err = os.Remove(f[0])

	if err != nil {
		services.Critical(errors.New(err.Error() + "Could not delete file - " + f[0]))
	}

	// Delete CSV Zip file
	err = os.Remove(zipFile)

	if err != nil {
		services.Critical(errors.New(err.Error() + "Could not delete file - " + zipFile))
	}

	// Return happy
	return options, underlyingLast, err
}

//
// downloadEodSymbol - Download one symbol and store it locally for back testing.
//
func downloadEodSymbol(symbol string, date time.Time) (string, error) {
	// Download File
	dFile := "options-eod/" + strings.ToUpper(symbol) + "/" + date.Format("2006-01-02") + ".csv.zip"
	lFile := os.Getenv("CACHE_DIR") + "/object-store/" + dFile

	// List file just to make sure we have this file at the storage
	list, err := object.ListObjects(dFile)

	if err != nil {
		return "", err
	}

	// We should return one file only.
	if len(list) != 1 {
		return "", errors.New("File not found at Object Store : " + dFile)
	}

	// Check the MD5 the current file. If we already have the file no need to re-download.
	md5Hash := files.Md5(lFile)

	// Send download job to the workers
	if md5Hash != strings.Replace(list[0].ETag, `"`, "", -1) {

		// Download file from object store.
		_, err2 := object.DownloadObject(dFile)

		if err2 != nil {
			return "", err2
		}
	}

	// Return happy
	return lFile, nil
}

//
// DownloadAllEodSymbolFiles - Download one symbol and store it locally for back testing.
//
// TODO(spicer): WE do not use this anywhere but some day we might to prime symbols.
// But most likely we will have to add the storeSymbolEodJSON (maybe just by calling GetOptionsBySymbol)
// to this no point in storing the CSV when JSON is what we really want.
//
func DownloadAllEodSymbolFiles(symbol string, debug bool) []string {

	var total int = 0
	var allFiles []string

	// Log if we are debug
	if debug {
		fmt.Println("Starting download of all " + symbol + " daily options data.")
	}

	// Worker job queue (some stocks have thousands of days)
	jobs := make(chan Job, 1000000)
	results := make(chan Job, 1000000)

	// Load up the workers
	for w := 0; w < workerCount; w++ {
		go downloadWorker(jobs, results)
	}

	// List files we need to download.
	list, err := object.ListObjects("options-eod/" + strings.ToUpper(symbol) + "/")

	if err != nil {
		services.Info(err)
		return allFiles
	}

	// Send all files to workers.
	for key, row := range list {

		// Check the MD5 the current file. If we already have the file no need to re-download.
		md5Hash := files.Md5(os.Getenv("CACHE_DIR") + "/object-store/" + row.Key)

		// Files we return.
		allFiles = append(allFiles, os.Getenv("CACHE_DIR")+"/object-store/"+row.Key)

		// Send download job to the workers
		if md5Hash != strings.Replace(row.ETag, `"`, "", -1) {
			total++
			jobs <- Job{Path: row.Key, Index: key}
		}
	}

	// Close jobs so the workers return.
	close(jobs)

	// Collect results so this function does not just return.
	for a := 0; a < total; a++ {
		job := <-results

		if debug {
			fmt.Println(job.Index, " of ", total)
		}
	}

	// Close results
	close(results)

	// Log if we are cli
	if debug {
		fmt.Println("Done download of all " + symbol + " daily options data.")
	}

	// Return a list of all files.
	return allFiles
}

//
// downloadWorker - A worker for downloading
//
func downloadWorker(jobs <-chan Job, results chan<- Job) {

	// Wait for jobs to come in and process them.
	for job := range jobs {

		// Download file from object store.
		_, err := object.DownloadObject(job.Path)

		if err != nil {
			services.Info(err)
			return
		}

		// Send back a happy result.
		results <- job

	}

}

//
// getSymbolEodJSONFromCache - See if we have the JSON for this symbol stored locally
//
func getSymbolEodJSONFromCache(symb string, date time.Time) ([]types.OptionsChainItem, float64, error) {
	// Read JSON file
	jdat, err := ioutil.ReadFile(cacheDir + symb + "/" + date.Format("2006-01-02") + ".json")

	if err != nil {
		return []types.OptionsChainItem{}, 0.00, err
	}

	// Convert json to struct
	var s SymbolStore
	if err := json.Unmarshal(jdat, &s); err != nil {
		return []types.OptionsChainItem{}, 0.00, err
	}

	// Return happy
	return s.Options, s.Last, nil
}

//
// storeSymbolEodJson - Take an options chain for a symbol and write it to our cache
//
func storeSymbolEodJSON(symb string, date time.Time, underlyingLast float64, options []types.OptionsChainItem) error {
	// Make sure our cache dir and symbol is setup.
	if _, err := os.Stat(cacheDir + symb); os.IsNotExist(err) {
		os.MkdirAll(cacheDir+symb, 0755)
	}

	// Create JSON for cache file
	jBlob, err := json.Marshal(SymbolStore{Options: options, Last: underlyingLast})

	if err != nil {
		return err
	}

	// Write the json to cache
	jFile := cacheDir + symb + "/" + date.Format("2006-01-02") + ".json"
	err = ioutil.WriteFile(jFile, jBlob, 0644)

	if err != nil {
		return errors.New(err.Error() + " : Could not write json file - " + jFile)
	}

	// Return happy
	return nil
}

/* End File */
