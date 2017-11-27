//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// This is a wrapper class for our object store or choice. (ie. AWS S3)

package object

import (
	"io"
	"os"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

//
// List files at object store.
//
func ListObjects(prefix string) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4)
	// is automatically determined based on the Endpoint value.
	s3Client, err := minio.New(os.Getenv("OBJECT_ENDPOINT"), os.Getenv("OBJECT_ACCESS_KEY_ID"), os.Getenv("OBJECT_SECRET_ACCESS_KEY"), true)

	if err != nil {
		return objects, err
	}

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	// Returns: https://github.com/minio/minio-go/blob/368f6da0bdf6d5bb6fbcbf6afdb626e711bab6df/api-datatypes.go#L34
	for object := range s3Client.ListObjects(os.Getenv("OBJECT_BUCKET"), prefix, true, doneCh) {

		if object.Err != nil {
			return objects, err
		}
		objects = append(objects, object)

	}

	// Return a happy array of objects
	return objects, nil
}

//
// Download an object to our cache directory.
//
func DownloadObject(objectPath string) (string, error) {

	// Set the cache dir.
	cacheDir := os.Getenv("CACHE_DIR") + "/object-store/"

	// Make a directory to download.
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0755)
	}

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4)
	// is automatically determined based on the Endpoint value.
	s3Client, err := minio.New(os.Getenv("OBJECT_ENDPOINT"), os.Getenv("OBJECT_ACCESS_KEY_ID"), os.Getenv("OBJECT_SECRET_ACCESS_KEY"), true)

	if err != nil {
		return "", err
	}

	object, err := s3Client.GetObject(os.Getenv("OBJECT_BUCKET"), objectPath, minio.GetObjectOptions{})

	if err != nil {
		return "", err
	}

	// Copy file to local local location.
	localFile, err := os.Create(cacheDir + filepath.Base(objectPath))

	if err != nil {
		return "", err
	}

	if _, err = io.Copy(localFile, object); err != nil {
		return "", err
	}

	// Return happy.
	return cacheDir + filepath.Base(objectPath), nil
}

/* End File */
