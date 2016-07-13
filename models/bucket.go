package models

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type MetaBucket struct {
	Name      string `json:"name"`
	TimeStamp int    `json:"timestamp"`

	// TODO: more metadata
}

type Bucket struct {
	Id   []byte     `json:"name"`
	Meta MetaBucket `json:"meta"`
}

type Buckets []Bucket
type BucketModel struct{}

func GetMeta(b *bolt.Bucket) (meta MetaBucket) {
	m := b.Get([]byte("meta"))
	json.Unmarshal(m, &meta)
	return meta
}

func (m BucketModel) One(bucketId []byte) (bucket Bucket) {
	db.View(func(tx *bolt.Tx) error {
		bucket = tx.Bucket(bucketId)
	})
	return bucket
}

func (m BucketModel) All() (buckets Buckets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bucketId []byte, b *bolt.Bucket) error {
			m := GetMeta(b)
			buckets = append(buckets, models.Bucket{Id: bucketId, Meta: m})
			return nil
		})
		return nil
	})
	return buckets, err
}
