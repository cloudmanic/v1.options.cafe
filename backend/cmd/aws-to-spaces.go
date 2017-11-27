//
// Date: 11/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	env "github.com/jpfuentes2/go-env"
)

//
// Main.
//
func main() {

	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load ENV (if we have it.)
	env.ReadEnv("../.env")

	ListFilesAWS()
}

//
// List files from AWS.
//
func ListFilesAWS() {

	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	if err != nil {
		panic(err)
	}

	input := &s3.ListObjectsInput{Bucket: aws.String(os.Getenv("AWS_BUCKET"))}

	result, err := svc.ListObjects(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}

/* End File */
