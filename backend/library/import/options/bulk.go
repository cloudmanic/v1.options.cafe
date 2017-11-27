//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package options

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/app.options.cafe/backend/library/files"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/app.options.cafe/backend/library/store/object"
	"github.com/araddon/dateparse"
)

//
// Here we do a bulk import of all End of Day options data.
// First we go over historical data stored in a per month zip file.
// Then we go through all the data we stored per day after historical data.
// For production typically we only run this one. As we update this data
// via cron every day. But this is useful for getting setup locally or
// rebuilding our production for some reason. We take EOD options data
// that we purchase from a service and then break it up into a file per day
// per symbol. We then store this file in an object store such as S3.
//
func DoBulkEodImportToPerSymbolDay() {

	// Do historical one month EOD data
	err := DoHistoricalMonthlyEodData()

	if err != nil {
		panic(err)
	}

	// Success processing.
	services.Info("DoHistoricalMonthlyEodData Done.")
}

//
// Do Historical Monthly EOD data processing.
//
func DoHistoricalMonthlyEodData() error {

	// Get a map of monthly EOD's that we have already done.
	monthlyDone, err := ReadMonthlyEodImportLog()

	if err != nil {
		return err
	}

	// Grab a list all historical (per month) data from our object store. (stored in bucket options-eod-monthly)
	// This goes from options-eod-monthly/2002_April.zip - 2017_September.zip
	objects, err := object.ListObjects("options-eod-monthly")

	if err != nil {
		return err
	}

	// Loop through each month / year zip file in options-eod-monthly.
	for _, row := range objects {

		// Check to see if we already did this key (maybe the program died and we are resuming)
		if _, ok := monthlyDone[filepath.Base(row.Key)]; ok {
			services.Info("Skipping (already done): " + row.Key)
			continue
		}

		// TODO: pass this off to a worker.
		// Process one Month's worth of EOD data.
		err := DoOneMonthsEodData(row.Key)

		if err != nil {
			return err
		}

	}

	// Return Happy.
	return nil
}

//
// Process One month's worth of EOD data.
//
func DoOneMonthsEodData(key string) error {

	// Downalod object
	path, err := DownloadOneMonthEodData(key)

	if err != nil {
		return err
	}

	// Process the data we just imported.
	err = OneMonthEodImport(path)

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

//
// Download one month's worth of EOD data.
//
func DownloadOneMonthEodData(key string) (string, error) {

	// Download file to our cache dir.
	services.Info("Downloading: " + key)

	path, err := object.DownloadObject(key)

	if err != nil {
		return "", err
	}

	// Return Happy.
	return path, nil
}

//
// Process one Month's worth of data.
//
func OneMonthEodImport(zipFilePath string) error {

	// Unzip CSV files.
	files, err := files.Unzip(zipFilePath, os.Getenv("CACHE_DIR"))

	if err != nil {
		return err
	}

	// Loop through the files paths and find the file we are looking for.
	for _, row := range files {

		// We are only interested in one file.
		i := strings.Index(row, "ptions_")

		if i <= -1 {
			err := os.Remove(row)

			if err != nil {
				return err
			}

			continue
		}

		// Make sure we have a file
		if len(row) <= 0 {
			continue
		}

		// Import the CSV file.
		err = OneDayEodImport(row)

		if err != nil {
			return err
		}

	}

	// Log this file is complete
	f, err := os.OpenFile(os.Getenv("CACHE_DIR")+"/OneMonthEodImport.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(filepath.Base(zipFilePath) + "\n"); err != nil {
		return err
	}

	// Return Happy
	return nil
}

//
// Import one day's worth of EOD options data. We pass in a path to a CSV File.
// We break this file up into one file per symbol. Then upload it to our data store.
//
func OneDayEodImport(csvFile string) error {

	// Open CSV file
	f, err := os.Open(csvFile)

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
	services.Info("Importing option EOD quotes for - " + string(lines[1][7]))

	// Figure out quote date
	date, err := dateparse.ParseAny(string(lines[1][7]))

	if err != nil {
		return err
	}

	// Hash map to keep each symbol array
	symbolMap := make(map[string][][]string)

	// Loop through lines & turn into object
	for _, line := range lines {

		// Skip heading line
		if line[0] == "UnderlyingSymbol" {
			continue
		}

		// Add to map.
		symbolMap[line[0]] = append(symbolMap[line[0]], line)
	}

	// Loop through and store to file.
	for key := range symbolMap {

		// Store Symbol to a file based on date and symbol
		_, err := OneDayEodSymbol(key, date, symbolMap[key])

		if err != nil {
			return err
		}

		//fmt.Println(zipFilePath)

		// // Send the file to AWS for storage
		// err = AWSUpload(zipFilePath, key)

		// if err != nil {
		// 	return err
		// }

	}

	// Delete output file.
	err = os.Remove(csvFile)

	if err != nil {
		return errors.New("Could not delete file - " + csvFile)
	}

	// Return happy.
	return nil
}

//
// Store one day's worth of EOD options Symbols
//
func OneDayEodSymbol(symbol string, date time.Time, data [][]string) (string, error) {

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
// Read Monthly log.
//
func ReadMonthlyEodImportLog() (map[string]bool, error) {
	log := make(map[string]bool)

	file, err := os.Open(os.Getenv("CACHE_DIR") + "/OneMonthEodImport.log")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		log[scanner.Text()] = true
	}

	return log, scanner.Err()
}

/* End File */
