//
// Date: 11/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	minio "github.com/minio/minio-go"
)

func worker(jobs <-chan string, results chan<- string) {

	accessKey := "DTCYU27W6JSXEHQGOCRS"
	secKey := "S1V7MrJI48F5av99pPKlhUyxhk3vQULsSYlO+rxkCVU"
	ssl := true

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New("nyc3.digitaloceanspaces.com", accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}

	for j := range jobs {

		object, err := client.GetObject("app-options-cafe", j, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}

		localFile, err := os.Create("/tmp/" + j)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err = io.Copy(localFile, object); err != nil {
			fmt.Println(err)
			return
		}

		results <- j
	}
}

//
// Main ...
//
func main() {

	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	jobs := make(chan string, 10000)
	results := make(chan string, 10000)

	for w := 1; w <= 100; w++ {
		go worker(jobs, results)
	}

	accessKey := "DTCYU27W6JSXEHQGOCRS"
	secKey := "S1V7MrJI48F5av99pPKlhUyxhk3vQULsSYlO+rxkCVU"
	ssl := true

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New("nyc3.digitaloceanspaces.com", accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	objectCh := client.ListObjects("app-options-cafe", "options-daily-eod/", false, doneCh)

	count := 0
	for obj := range objectCh {
		if obj.Err != nil {
			fmt.Println(obj.Err)
			return
		}

		jobs <- obj.Key
		count++

		// fmt.Println(obj.Key)

		// object, err := client.GetObject("app-options-cafe", obj.Key, minio.GetObjectOptions{})
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// localFile, err := os.Create("/tmp/" + obj.Key)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// if _, err = io.Copy(localFile, object); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

	}
	close(jobs)

	for a := 1; a <= count; a++ {
		key := <-results
		fmt.Println(key)
	}
}
