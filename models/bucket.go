package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

// Bucket struct.
// Used for sending as json.
type SnippetBucket struct {
	Name string `json:"name"`
}
type SnippetBuckets []SnippetBucket

func WriteBucketToFileSys(storageRootPath string, bucketname string, tx *bolt.Tx) (err error) {

	bucketRootPath := storageRootPath + bucketname + "/"

	var snippets []Snippet
	snippets, err = IndexSnippets(bucketname, tx)
	if err != nil {
		fmt.Printf("Error indexing snippets: %v", err)
	}
	for _, snippet := range snippets {
		cleanFullName := filepath.Clean(snippet.Name)
		fullFilePath := filepath.Dir(cleanFullName)
		if fullFilePath == "." {
			fullFilePath = ""
		}
		err = os.MkdirAll(bucketRootPath+fullFilePath, 0777)                                                                //rw
		err = ioutil.WriteFile(bucketRootPath+fullFilePath+"/"+filepath.Base(cleanFullName), []byte(snippet.Content), 0666) //rw, truncates before write
	}
	return err
}
