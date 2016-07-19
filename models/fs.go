package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"

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

var snippetModel = new(SnippetModel)
var bucketModel = new(BucketModel)

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
	bucketName = folders[0]
	fmt.Printf("Got bucketName: %v\n:", bucketName)
	return bucketName
}

func getSnippetNamebyFSPath(path string) (snippetName string) {
	c := filepath.Clean(path)
	d := strings.Replace(c, config.FSStorePath+"/", "", 1) // remove hacks/
	b := getBucketNameByFSPath(path)                       //
	snippetName = strings.Replace(d, b, "", 1)             // remove bucket/
	fmt.Printf("Got snippetName: %v\n:", snippetName)
	return snippetName
}

func filefySnippet(bucket Bucket, snippetId string) (err error) {

	var snippet Snippet

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.Id))
		s := b.Get([]byte(snippetId))
		e := json.Unmarshal(s, &snippet)
		return e
	})
	if err != nil {
		return err
	}

	// Write file with all necessary dir paths to it.
	cleanName := filepath.Clean(snippet.Name)
	d := filepath.Join(config.FSStorePath, bucket.Meta.Name, filepath.Dir(cleanName))
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return err
	}

	// Write contents of file (truncates).
	err = ioutil.WriteFile(filepath.Join(d, filepath.Base(cleanName)), []byte(snippet.Content), 0666)
	return err
}

func (fs FSModel) WriteFile(bucketId string, snippetId string) (err error) {

	// Get bucket and snippet.
	var bucket Bucket

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketId))
		bucket.Id = bucketId
		bucket.Meta = getMeta(b)

		return nil
	})

	return filefySnippet(bucket, snippetId)
}

func (fs FSModel) WriteDirHard(bucketId string) (err error) {

	var bucket Bucket

	err = db.View(func(tx *bolt.Tx) error {

		// Setup bucket we're working on.
		b := tx.Bucket([]byte(bucketId))
		bucket.Id = bucketId
		bucket.Meta = getMeta(b)

		// Remove all --> hard.
		rmdirerr := os.RemoveAll(filepath.Join(config.FSStorePath, bucket.Meta.Name))
		if rmdirerr != nil {
			return rmdirerr
		}

		// Iterate through snippets, writing a file for each.
		berr := b.ForEach(func(k, v []byte) error {
			if string(k) != "meta" {
				return filefySnippet(bucket, string(k))
			} else {
				return nil
			}
		})
		return berr // should be nil
	})
	return err // should be nil
}

func bucketizeDir(path string) (bucket Bucket, err error) {
	// Extract bucket name from path.
	bucketName := getBucketNameByFSPath(path)
	// Check if bucket exists by name.
	bucket, err = GetBucketByName(bucketName) // Bucket{}
	if err != nil {
		return bucket, err
	}
	fmt.Println("Got bucket: %v\n", bucket)

	// If bucket doesn't exist by name (from hacks/@bucket/path/to/file)
	if bucket == (Bucket{}) {
		// Initialize a new bucket.
		fmt.Println("Creating new bucket...")
		bucket, err = bucketModel.Create(bucketName)
		fmt.Println("New bucket: %v\n", bucket)
		if err != nil {
			return bucket, err
		}
	}
	return bucket, err
}

func snippetizeFile(path string) (bucket Bucket, snippet Snippet, err error) {

	fmt.Printf("Snippetizering path: %v\n", path)

	// Read contents of file.
	contents, ioerr := ioutil.ReadFile(path) // func ReadFile(filename string) ([]byte, error)
	if ioerr != nil {
		fmt.Println("IO err: %v\n", ioerr)
		return bucket, snippet, err
	}
	fmt.Println("File contents: %v\n", contents)

	bucket, err = bucketizeDir(path)
	if err != nil {
		return bucket, snippet, err
	}
	fmt.Println("Bucket id: %v\n", bucket.Id)

	// Extract snippet name from path.
	snippetName := getSnippetNamebyFSPath(path)
	fmt.Println("Snippet name: %v\n", snippet.Name)

	// Check if snippet exists by name.
	snippet = GetSnippetByName(bucket.Id, snippetName)

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

	fmt.Println("Snippet: %v\n", snippet)

	fmt.Println("Setting snippet...")
	err = snippetModel.Set(snippet)

	return bucket, snippet, err
}

// Note: accepts FULL path (includes FSStoreDir)
func (fs FSModel) SnippetizeFile(path string) (bucket Bucket, snippet Snippet, err error) {
	return snippetizeFile(path)
}

func (fs FSModel) SnippetizeDir(path string) (buckets Buckets, snippets Snippets, err error) {
	err = filepath.Walk(filepath.Clean(path), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !(info.Name() == ".stfolder") {
			bucket, snippet, snipErr := snippetizeFile(path)
			if snipErr != nil {
				// TODO: Handle errors better. Like in their own model.
			}
			buckets = append(buckets, bucket)
			snippets = append(snippets, snippet)
		}
		return nil
	})

	return buckets, snippets, err
}

func DeleteFile(path string) error {
	p := filepath.Clean(path)
	return os.Remove(p)
}

func DeleteDir(path string) error {
	p := filepath.Clean(path)
	return os.RemoveAll(p)
}
