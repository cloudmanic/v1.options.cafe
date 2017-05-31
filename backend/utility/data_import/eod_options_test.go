//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import (
  "os"
  //"fmt"
  "testing"  
  "github.com/joho/godotenv"  
  "github.com/tj/go-dropy"
  "github.com/tj/go-dropbox"
  "github.com/stretchr/testify/assert"  
)

func TestGetProccessedFiles(t *testing.T) {

  // Load .env file 
  err := godotenv.Load("../../.env")
  assert.NoError(t, err, "Error loading .env file")

  // Get the Dropbox Client (this is where we archive the zip file)
  client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
  
  // Get files we should skip.
  skip, err := GetProccessedFiles(client)
  assert.NoError(t, err)
  assert.NotEqual(t, 0, len(skip), "Should have more than one skipped file")  

}

/* End File */