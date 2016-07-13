package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

type Snippet struct {
	Id        string `json:"id"`
	BucketId  string `json:"bucketId"`
	Name      string `json:"name"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	TimeStamp int    `json:"timestamp"`
	Meta      string `json:"meta"`
}
type Snippets []Snippet
type SnippetModel struct{}

func snipFromJSON(snippetJSONBytes []byte) (snippet Snippet) {
	json.Unmarshal(snippetJSONBytes, &snippet)
	return snippet
}

func GetSnippetByName(bucketname string, name string, tx *bolt.Tx) (snippet Snippet) {
	b := tx.Bucket([]byte(bucketname))
	c := b.Cursor()

	for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
		snip := snipFromJSON(snipval)
		if snip.Name == name {
			snippet = snip
			break
		}
	}
	if snippet == (Snippet{}) {
		return Snippet{}
	} else {
		return snippet
	}
}

func (m SnippetModel) All(bucketId string) (snippets Snippets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketId))
		c := b.Cursor()
		for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
			var snip Snippet
			json.Unmarshal(snipval, &snip)
			snippets = append(snippets, snip)
		}
		return nil
	})
	return snippets, err
}

func (m SnippetModel) Set(snippet Snippet) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(snippet.BucketId))
		j, _ := json.Marshal(snippet)
		return b.Put([]byte(snippet.Id), j)
	})
}

func (m SnippetModel) Delete(bucketid string, snippetid string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketid))
		// Get snippet so can access Name (to remove from FS).
		var snip Snippet
		v := b.Get([]byte(snippetid))
		json.Unmarshal(v, &snip)

		// First, remove from bucket.
		derr := b.Delete([]byte(snip.Id))
		if derr != nil {
			fmt.Printf("Error deleting from bucket: %v", derr)
			return derr
		} else {
			// Remove from FS if was successfully deleted from bucket.
			path := "./hacks/snippets/" + snip.Name
			fmt.Printf("Removing file at path: %v", path)
			derr = os.Remove(path)
			if derr != nil {
				fmt.Printf("Error removing file: %v", derr)
			}
		}
		return derr
	})
	return err
}
