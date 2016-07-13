package models

import (
	"encoding/json"
	"fmt"
	"time"

	"../config"
	"github.com/boltdb/bolt"
)

// Set global var.
// This will the pointer to our open BoltDB.
var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open(config.BoltDBPath, 0666, nil)
	if err != nil {
		fmt.Printf("Could not initialize Bolt database. Error: %v\n", err)
	}

	// Ensure existence of default bucket.
	db.Update(func(tx *bolt.Tx) error {
		var defaultBucket Bucket
		var defaultMetaBucket MetaBucket

		defaultMetaBucket.Name = config.DefaultBucketName
		defaultMetaBucket.TimeStamp = int(time.Now().UTC().Unix() * 1000)

		defaultBucket.Meta = defaultMetaBucket

		j, _ := json.Marshal(defaultMetaBucket)

		b, aerr := tx.CreateBucketIfNotExists([]byte(config.DefaultBucketName))
		if aerr != nil {
			fmt.Errorf("create bucket err: %s", aerr)
		} else {
			fmt.Printf("create bucket: %v", b)

			// Create meta info if it doesn't exist.
			bb := b.Get([]byte("meta"))
			if bb == nil {
				b.Put([]byte("meta"), j)
			}
		}
		return aerr
	})
}

// Getter.
func GetDB() *bolt.DB {
	return db
}
