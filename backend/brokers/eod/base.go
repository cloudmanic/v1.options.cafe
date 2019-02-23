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
	"fmt"
	"go/build"
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/files"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/store/object"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
	minio "github.com/minio/minio-go"
)

const workerCount int = 100

type Api struct {
	DB  models.Datastore
	Day time.Time // This is the day we pull EOD data for
}

type Job struct {
	Path  string
	Index int
}

//
// Start up the controller.
//
func init() {
	// Helpful for testing
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")
}

//
// Get options by Symbol. We do not return a chain. We more or less
// return the data in our CSV files as Go structs
//
func (t *Api) GetOptionsBySymbol(symbol string) ([]types.OptionsChainItem, float64, error) {

	symb := strings.ToUpper(symbol)
	underlyingLast := 0.00
	options := []types.OptionsChainItem{}

	// Set the cache dir.
	cacheDir := os.Getenv("CACHE_DIR") + "/object-store/options-eod/"

	// Get dates from S3 store - TODO: Maybe some day just download the date we are after instead of all dates.
	DownloadEodSymbol(symbol, false)

	// Make sure we have this zip file
	zipFile := cacheDir + symb + "/" + t.Day.Format("2006-01-02") + ".csv.zip"

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
			Underlying:     symb,
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

	// // Delete csv file
	err = os.Remove(f[0])

	if err != nil {
		services.FatalMsg(err, "Could not delete file - "+f[0])
	}

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

//
// Download one symbol and store it locally for back testing.
//
func DownloadEodSymbol(symbol string, debug bool) []string {

	var total int = 0
	var allFiles []string

	// Cache keys
	cacheKey := "oc-brokers-eod-symbol-objects-" + strings.ToUpper(symbol)
	cacheListKey := "oc-brokers-eod-local-cache-list-" + strings.ToUpper(symbol)

	// See if we have this result in the cache.
	var cachedFileList []string
	found1, _ := cache.Get(cacheListKey, &cachedFileList)

	// Store cache
	if found1 {
		return cachedFileList
	}

	// Log if we are debug
	if debug {
		fmt.Println("Starting download of all " + symbol + " daily options data.")
	}

	// Worker job queue (some stocks have thousands of days)
	jobs := make(chan Job, 1000000)
	results := make(chan Job, 1000000)

	// Load up the workers
	for w := 0; w < workerCount; w++ {
		go DownloadWorker(jobs, results)
	}

	var list []minio.ObjectInfo

	// See if we have this result in the cache.
	cachedObjectInfo := []minio.ObjectInfo{}
	found, _ := cache.Get(cacheKey, &cachedObjectInfo)

	// Store cache
	if found {

		list = cachedObjectInfo

	} else {

		// List files we need to download.
		list, err := object.ListObjects("options-eod/" + strings.ToUpper(symbol) + "/")

		if err != nil {
			services.Warning(err)
			return allFiles
		}

		// Store dates in cache - 3 hours
		cache.SetExpire(cacheKey, (time.Minute * 60 * 3), list)

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

	// Log if we are cli
	if debug {
		fmt.Println("Done download of all " + symbol + " daily options data.")
	}

	// Store dates in cache - 3 hours
	cache.SetExpire(cacheListKey, (time.Minute * 60 * 3), allFiles)

	// Return a list of all files.
	return allFiles
}

//
// A worker for downloading
//
func DownloadWorker(jobs <-chan Job, results chan<- Job) {

	// Wait for jobs to come in and process them.
	for job := range jobs {

		// Download file from object store.
		_, err := object.DownloadObject(job.Path)

		if err != nil {
			services.Warning(err)
			return
		}

		// Send back a happy result.
		results <- job

	}

}

/* End File */
