//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package options

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/cloudmanic/app.options.cafe/backend/library/files"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/store/object"
	"github.com/jlaffaye/ftp"
	minio "github.com/minio/minio-go"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

//
// Do End of Day Options Import. We run this every day from "Cron". It will
// connect to the FTP site at DeltaNeutral (our EOD options data provider) and
// download any days data that we do not already have. Once we download the file we
// upload it to Dropbox for achieving. Then we break up the DeltaNeutral export
// into one file per date per asset symbol. Lastly we upload this data to AWS.
//
func DoEodOptionsImport() {

	// Log
	services.Info("Starting DoEodOptionsImport().")

	// Get files we should skip.
	proccessedFiles, err := GetProccessedFiles()

	if err != nil {
		services.FatalMsg(err, "Could not get proccessedFiles files in GetProccessedFiles()")
		return
	}

	// Get daily options from vendor
	ftpFiles, err := GetOptionsDailyData()

	if err != nil {
		services.FatalMsg(err, "Could not get ftp files in GetOptionsDailyData()")
		return
	}

	// Loop through the FTP files.
	for _, row := range ftpFiles {

		// Have we already processed this day?
		if _, ok := proccessedFiles[row.Name]; ok {
			continue
		}

		// Processing Log
		services.Info("Processing " + row.Name + " from data provider.")

		// Download file.
		filePath, err := DownloadOptionsDailyDataByName(row.Name)

		if err != nil {
			services.FatalMsg(err, "Could not download ftp file in DownloadOptionsDailyDataByName() - "+row.Name)
			continue
		}

		// Import file into different symbols (we do this first and then upload the archive we if things go wrong we can re run this.)
		err = SymbolImport(filePath)

		if err != nil {
			services.FatalMsg(err, "Could not import DatabaseImport() - "+row.Name)
			continue
		}

		// Upload file to object store.
		err = object.UploadObject(filePath, "options-eod-daily/"+filepath.Base(filePath))

		if err != nil {
			services.FatalMsg(err, "Could not upload to Object Store - "+row.Name)
			continue
		}

		// Log file
		services.Info("Finished uploading " + row.Name + " to Object Store.")

		// Open file we upload.
		file, err := os.Open(filePath)

		if err != nil {
			services.FatalMsg(err, "Could not open ftp file os.Open() - "+row.Name)
			continue
		}

		// Upload the file to Dropbox (just for safe keeping, all processing is done with the object store.)
		client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
		client.Upload("/data/options-eod-daily/"+row.Name, file)

		if err != nil {
			services.FatalMsg(err, "Could not upload to Dropbox - "+row.Name)
			continue
		}

		// Log file
		services.Info("Finished uploading " + row.Name + " to Dropbox.")

		// // Delete file we uploaded to Dropbox
		err = os.Remove(filePath)

		if err != nil {
			services.FatalMsg(err, "Could not delete file - "+filePath)
			continue
		}

	}

	// Send health check notice.
	if len(os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL")) > 0 {

		resp, err := http.Get(os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL"))

		if err != nil {
			services.FatalMsg(err, "Could send health check - "+os.Getenv("HEALTH_CHECK_DOEODOPTIONSIMPORT_URL"))
		}

		defer resp.Body.Close()

	}

	// Log
	services.Info("Done DoEodOptionsImport().")
}

//
// Break the data up per symbol
//
func SymbolImport(filePath string) error {

	// Log
	services.Info("Start SymbolImport - " + filePath)

	// Unzip CSV files.
	files, err := files.Unzip(filePath, "/tmp/output/")

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
		} else {
			err := os.Remove(row)

			if err != nil {
				services.FatalMsg(err, "Could not delete file - "+row)
			}
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
	services.Info("Importing option EOD quotes for - " + string(lines[0][7]))

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
	for key := range symbolMap {

		// Store Symbol to a file based on date and symbol
		zipFile, err := StoreOneDaySymbol(key, date, symbolMap[key])

		if err != nil {
			return err
		}

		// Upload to the object store.
		err = object.UploadObject(zipFile, "options-eod/"+key+"/"+filepath.Base(zipFile))

		if err != nil {
			return errors.New("object.UploadObject: " + err.Error() + " - " + key + " - " + zipFile)
		}

		// Delete zipFile file.
		err = os.Remove(zipFile)

		if err != nil {
			return errors.New("Could not delete file - " + zipFile)
		}

	}

	// Delete output file.
	err = os.Remove(file)

	if err != nil {
		services.FatalMsg(err, "Could not delete file - "+file)
	}

	// Log
	services.Info("Done SymbolImport - " + filePath)

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
	csvFilePtr, err := os.Create(dirPath + fileName)

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
	err = files.ZipFiles(zipFile, []string{csvFile})

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
func GetProccessedFiles() (map[string]minio.ObjectInfo, error) {

	// Files that we already have imported
	skip := make(map[string]minio.ObjectInfo)

	// Get all files stored at AWS object store.
	files, err := object.ListObjects("options-eod-daily")

	if err != nil {
		return nil, err
	}

	// Create a map (hash table) of all the files we already have at Object store
	for _, row := range files {
		skip[filepath.Base(row.Key)] = row
	}

	// Return happy.
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
