package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

type Snippet struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	TimeStamp int    `json:"timestamp"`
	Meta      string `json:"meta"`
}
type Snippets []Snippet

func SnipFromJSON(snippetBytes []byte) (snippet Snippet) {
	json.Unmarshal(snippetBytes, &snippet)
	return snippet
}

func IndexSnippets(bucketname string, tx *bolt.Tx) (snippets Snippets, err error) {
	b := tx.Bucket([]byte(bucketname))
	c := b.Cursor()

	for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
		var snip Snippet
		json.Unmarshal(snipval, &snip)
		snippets = append(snippets, snip)
	}

	return snippets, err
}

func SetSnippet(snippetid string, contents []byte, bucketname string, tx *bolt.Tx) (err error) {
	b := tx.Bucket([]byte(bucketname))
	err = b.Put([]byte(snippetid), contents)

	if err != nil {
		return fmt.Errorf("putting to bucket: %s", err)
	}
	return err
}

func DeleteSnippet(snippetid string, bucketid string, tx *bolt.Tx) (err error) {
	b := tx.Bucket([]byte(bucketid))

	// remove from os
	// get snippet so can access Name
	var snip Snippet
	v := b.Get([]byte(snippetid))
	json.Unmarshal(v, &snip)

	// remove from bucket
	err = b.Delete([]byte(snip.Id))
	if err != nil {
		fmt.Printf("Error deleting from bucket: %v", err)
		return err
	} else {
		// remove from bucket if was successfully deleted from bucket
		path := "./hacks/snippets/" + snip.Name
		fmt.Printf("Removing file at path: %v", path)
		removeErr := os.Remove(path)
		if removeErr != nil {
			fmt.Printf("Error removing file: %v", removeErr)
			// don't worry about not deleting a file that doesn't exist
		}
	}
	return err // return the error about deleting from bolt only
}
