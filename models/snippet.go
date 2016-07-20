package models

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"../config"
	"github.com/boltdb/bolt"
)

type Snippet struct {
	Id          string `json:"id"`
	BucketId    string `json:"bucketId"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	TimeStamp   int    `json:"timestamp"`
	Description string `json:"description"`
	IpCity      string `json:"ipCity"`
	Ip          string `json:"ip"`
}
type Snippets []Snippet
type SnippetModel struct{}

func snipFromJSON(snippetJSONBytes []byte) (snippet Snippet) {
	json.Unmarshal(snippetJSONBytes, &snippet)
	return snippet
}

func GetSnippetByName(bucketId string, name string) (snippet Snippet) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketId))
		c := b.Cursor()
		for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
			snip := snipFromJSON(snipval)
			if (snip.Name == name) || ("/"+snip.Name == name) || (snip.Name == "/"+name) {
				snippet = snip
				break
			}
		}
		return nil
	})

	if snippet == (Snippet{}) {
		return Snippet{}
	} else {
		return snippet
	}
}

func (m SnippetModel) UberAll() (snippets Snippets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketId []byte, b *bolt.Bucket) error {
			c := b.Cursor()
			for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
				snip := Snippet{}
				json.Unmarshal(snipval, &snip)

				if len(snip.Id) > 0 {
					snippets = append(snippets, snip)
				}

			}
			return nil
		})
	})
	return snippets, err
}

func (m SnippetModel) All(bucketId string) (snippets Snippets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketId))
		if b == nil {
			return nil
		}

		if b.Stats().KeyN > 0 {

			c := b.Cursor()
			for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
				var snip Snippet
				json.Unmarshal(snipval, &snip)
				snippets = append(snippets, snip)
			}
			return nil

		} else {
			return nil
		}
	})
	return snippets, err
}

func (m SnippetModel) Set(snippet Snippet) error {

	// fmt.Printf("Will try to set snip.\n")
	// fmt.Printf("snip.Id: %v\n", snippet.Id)
	// fmt.Printf("snippet.BucketId: %v\n", snippet.BucketId)
	// fmt.Printf("snippet.Name: %v\n", snippet.Name)
	// fmt.Printf("snippet.Language: %v\n", snippet.Language)
	// fmt.Printf("snippet.Content: %v\n", snippet.Content)
	// fmt.Printf("snippet.TimeStamp: %v\n", snippet.TimeStamp)
	// fmt.Printf("snippet.IpCity: %v\n", snippet.IpCity)

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(snippet.BucketId))
		// if berr != nil {
		// 	fmt.Printf("Could not create bucket if not exists for bucketId: %v\n", snippet.BucketId)
		// 	fmt.Printf("The error was: %v\n", berr)
		// 	return berr
		// }
		if b == nil {
			fmt.Printf("Could not find bucket with snippet.BucketId: %v\n", snippet.BucketId)
			return nil
		}
		j, err := json.Marshal(snippet)
		if err != nil {
			fmt.Printf("Could not marshal json: ")
			fmt.Printf("Error was : %v\n", err)
			return err
		}
		perr := b.Put([]byte(snippet.Id), j)
		if perr != nil {
			fmt.Printf("Error putting snippet to bucket: %v\n", perr)
			return perr
		}
		fmt.Printf("It would appear I successfully put hte snippet.")
		return nil
	})
}

func (m SnippetModel) Delete(bucketId string, snippetId string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketId))
		bm := getMeta(b)
		// Get snippet so can access Name (to remove from FS).
		var snip Snippet
		v := b.Get([]byte(snippetId))
		json.Unmarshal(v, &snip)

		// First, remove from bucket.
		derr := b.Delete([]byte(snip.Id))
		if derr != nil {
			fmt.Printf("Error deleting from bucket: %v", derr)
			return derr
		}

		// Remove file from FS.
		path := filepath.Join(config.FSStorePath, bm.Name, snip.Name)
		fderr := DeleteFile(path)
		if fderr != nil {
			return fderr
		}
		return nil
	})
	return err
}
