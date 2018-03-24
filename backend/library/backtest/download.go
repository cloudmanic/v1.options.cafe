//
// Date: 2018-03-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// Concurrent 1:   14m 12.401s
// Concurrent 10:  1m38.412s
// Concurrent 25:  0m56.556s
// Concurrent 50:  0m44.670s
// Concurrent 100: 0m27.545s
// Concurrent 200: Started to get errors at this point
//

package backtest

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/api.triwou.org/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/files"
	"github.com/cloudmanic/app.options.cafe/backend/library/store/object"
)

const workerCount int = 100

type Job struct {
	Path  string
	Index int
}

//
// Download one symbol and store it locally for back testing.
//
func DownloadEodSymbol(symbol string, cli bool) []string {

	var total int = 0
	var allFiles []string

	// Log if we are cli
	if cli {
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

		if cli {
			fmt.Println(job.Index, " of ", total)
		}
	}

	// Log if we are cli
	if cli {
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
