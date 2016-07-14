package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func (fs FSModel) WalkDir() (filepaths Filepaths, err error) {
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
