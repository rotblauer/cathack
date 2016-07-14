package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../config"
)

// TODO: is this being used anywhere?
// I don't think soo...
// type FS struct {
// 	Path string `json:"filepath"`
// }

type Filepath struct {
	Path     string   `json:"path"`
	FileInfo FileInfo `json:"fileInfo"` // ~~ os.FileInfo interface
}
type Filepaths []Filepath

type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
}

type FSModel struct{}

// I want to use FS to walk the HacksRootDir and collect all the
// files (and empty directories).

func (fs FSModel) CollectDirPaths() (filepaths Filepaths, err error) {
	err = filepath.Walk(config.FSStorePath, func(path string, info os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}

		f := FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		i := Filepath{
			Path:     path,
			FileInfo: f,
		}

		filepaths = append(filepaths, i)
		return nil
	})
	return filepaths, err
}

func getBucketNameByFSPath(path string) (bucketName string) {
	c := filepath.Clean(path)
	d := filepath.Dir(c)
	withinHacksRootDir := strings.Replace(dir, config.FSStorePath+"/", "", 1)
	folders := strings.Split(withinHacksRootDir, "/")
	return folders[0]
}

func getSnippetNamebyFSPath(path string) (snippetName string) {
	c := filepath.Clean(path)
	d := strings.Replace(c, config.FSStorePath+"/", "", 1) // remove hacks/
	b := getBucketNameByFSPath(path)                       //
	snippetName = strings.Replace(d, b, "", 1)             // remove bucket/
}

// TODO: route this.?
func (fs FSModel) SetFile(bucket Bucket, snippet Snippet) (err error) {
	cleanName := filepath.Clean(snippet.Name)
	d := filepath.Join(config.FSStorePath, bucket.Meta.Name, filepath.Dir(cleanName))
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(d, filepath.Base(cleanName)), []byte(snippet.Content), 0666)
	return err
}

// Note: accepts FULL path (includes FSStoreDir)
func (fs FSModel) SnippetizeFile(path string) (snippet Snippet, err error) {
	contents, ioerr := ioutil.ReadFile(path) // func ReadFile(filename string) ([]byte, error)
	if ioerr != nil {
		return nil, err
	}

	bucketName := getBucketNameByFSPath(path)
	bucket := findBucketByName(bucketName) // Bucket{}

	db.Update(func(tx *bolt.TX) error {
		b := tx.Bucket([]byte(bucket.Id))

	})
}

func (fs FSModel) SnippetizeDir(path string) (snippets Snippets) {

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
