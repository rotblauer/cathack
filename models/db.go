package models

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	dbLocation        = "hack.db"
	defaultBucketName = "snippets"
)

// Set global var.
// This will the pointer to our open BoltDB.
var db *bolt.DB

func init() {
	db, err = bolt.Open(dbLocation, 0666, nil)
	if err != nil {
		fmt.Printf("Could not initialize Bolt database. Error: %v\n", err)
	}

	// Ensure existence of default bucket.
	db.Update(func(tx *bolt.Tx) error {
		b, aerr := tx.CreateBucketIfNotExists([]byte(defaultBucketName))
		if aerr != nil {
			fmt.Errorf("create bucket err: %s", aerr)
		} else {
			fmt.Printf("create bucket: %v", b)
		}
		return aerr
	})
}

// Getter.
func GetDB() *bolt.DB {
	return db
}
