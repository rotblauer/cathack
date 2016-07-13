package models

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"../config"
)

type FS struct {
	Path string `json:"filepath"`
}

type FSModel struct{}

func (fs FSModel) SetFile(snippet Snippet) (err error) {
	cleanName := filepath.Clean(snippet.Name)
	d := filepath.Join(config.FSStorePath, snippet.BucketName, filepath.Dir(cleanName))
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(d, filepath.Base(cleanName)), []byte(snippet.Content), 0666)
	return err
}

// func (fs FSModel) SetDir() error {

// }

// func (fs FSModel) GetOne() error {

// }

// func (fs FSModel) GetAll() error {

// }

// func (fs FSModel) DeleteFile() error {

// }

func (fs FSModel) DeleteDir(path string) error {
	return os.RemoveAll(path)
}

// func WriteBucketToFileSys(storageRootPath string, bucketname string, tx *bolt.Tx) (err error) {

// 	bucketRootPath := storageRootPath + "/" + bucketname + "/"

// 	var snippets []Snippet
// 	snippets, err = IndexSnippets(bucketname, tx)
// 	if err != nil {
// 		fmt.Printf("Error indexing snippets: %v", err)
// 	}
// 	for _, snippet := range snippets {
// 		cleanFullName := filepath.Clean(snippet.Name)
// 		fullFilePath := filepath.Dir(cleanFullName)
// 		if fullFilePath == "." {
// 			fullFilePath = ""
// 		}
// 		err = os.MkdirAll(bucketRootPath+fullFilePath, 0777)                                                                //rw
// 		err = ioutil.WriteFile(bucketRootPath+fullFilePath+"/"+filepath.Base(cleanFullName), []byte(snippet.Content), 0666) //rw, truncates before write
// 	}
// 	return err
// }
