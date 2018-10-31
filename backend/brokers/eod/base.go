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
	"fmt"
	"go/build"
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/files"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/store/object"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
)

const workerCount int = 100

type Api struct {
	DB models.Datastore
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

	// List files we need to download.
	list, err := object.ListObjects("options-eod/" + strings.ToUpper(symbol) + "/")

	if err != nil {
		services.Warning(err)
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

	// Log if we are cli
	if debug {
		fmt.Println("Done download of all " + symbol + " daily options data.")
	}

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
