//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import (
	"bytes"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//
// Upload to S3
//
func AWSUpload(filePath string, symbol string) error {

	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	if err != nil {
		return err
	}

	// Upload
	err = AddFileToS3(s, filePath, symbol)

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

//
// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
//
func AddFileToS3(s *session.Session, fileDir string, symbol string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_BUCKET")),
		Key:           aws.String("options-eod/" + symbol + "/" + path.Base(fileDir)),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		//ServerSideEncryption: aws.String("AES256"),
	})

	return err
}

/* End File */
