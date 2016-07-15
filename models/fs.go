package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../config"
	"../lib"
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
	withinHacksRootDir := strings.Replace(d, config.FSStorePath+"/", "", 1)
	folders := strings.Split(withinHacksRootDir, "/")
	return folders[0]
}

func getSnippetNamebyFSPath(path string) (snippetName string) {
	c := filepath.Clean(path)
	d := strings.Replace(c, config.FSStorePath+"/", "", 1) // remove hacks/
	b := getBucketNameByFSPath(path)                       //
	snippetName = strings.Replace(d, b, "", 1)             // remove bucket/
	return snippetName
}

// TODO: route this.?
func (fs FSModel) WriteFile(bucket Bucket, snippet Snippet) (err error) {
	cleanName := filepath.Clean(snippet.Name)
	d := filepath.Join(config.FSStorePath, bucket.Meta.Name, filepath.Dir(cleanName))
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(d, filepath.Base(cleanName)), []byte(snippet.Content), 0666)
	return err
}

func bucketizeDir(path string) (bucket Bucket, err error) {
	// Extract bucket name from path.
	bucketName := getBucketNameByFSPath(path)
	// Check if bucket exists by name.
	bucket, err = findBucketByName(bucketName) // Bucket{}
	if err != nil {
		return bucket, err
	}
	// If bucket doesn't exist by name (from hacks/@bucket/path/to/file)
	if bucket == (Bucket{}) {
		// Initialize a new bucket.
		bucketModel := new(BucketModel)
		bucket, err = bucketModel.Create(bucketName)
		if err != nil {
			return bucket, err
		}
	}
	return bucket, err
}

func snippetizeFile(path string) (bucket Bucket, snippet Snippet, err error) {
	// Read contents of file.
	contents, ioerr := ioutil.ReadFile(path) // func ReadFile(filename string) ([]byte, error)
	if ioerr != nil {
		return bucket, snippet, err
	}

	bucket, err = bucketizeDir(path)
	if err != nil {
		return bucket, snippet, err
	}

	// Extract snippet name from path.
	snippetName := getSnippetNamebyFSPath(path)
	// Check if snippet exists by name.
	snippet = getSnippetByName(bucket.Id, snippetName)
	// If snippet doesn't exists by name (from (hacks/@bucketName/<this/is/a/filename.txt>)
	if snippet == (Snippet{}) {
		snippet.Name = snippetName
		snippet.BucketId = bucket.Id
		snippet.Id = lib.RandSeq(6)
	}
	// Set attributes for previously existing or fresh-off-the-press snippet.
	snippet.Content = string(contents)
	snippet.Language = lib.GetLanguageModeByExtension(snippetName)
	snippet.TimeStamp = int(time.Now().UTC().UnixNano() / 1000000)

	snippetModel := new(SnippetModel)
	err = snippetModel.Set(snippet)
	return bucket, snippet, err
}

// Note: accepts FULL path (includes FSStoreDir)
func (fs FSModel) SnippetizeFile(path string) (bucket Bucket, snippet Snippet, err error) {
	return snippetizeFile(path)
}

func (fs FSModel) SnippetizeDir(path string) (buckets Buckets, snippets Snippets, err error) {
	werr := filepath.Walk(config.FSStorePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			bucket, snippet, snipErr := snippetizeFile(path)
			if snipErr != nil {
				// TODO: Handle errors better. Like in their own model.
			}
			buckets = append(buckets, bucket)
			snippets = append(snippets, snippet)
		}
		return nil
	})
	return buckets, snippets, werr
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
