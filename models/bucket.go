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
	Id   []byte     `json:"id"`
	Meta MetaBucket `json:"meta"`
}

type Buckets []Bucket
type BucketModel struct{}

func getMeta(b *bolt.Bucket) (meta MetaBucket) {
	m := b.Get([]byte("meta"))
	json.Unmarshal(m, &meta)
	return meta
}

func (m BucketModel) One(bucketId []byte) (bucket Bucket) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketId)
		bucket.Id = bucketId
		bucket.Meta = getMeta(b)
		return nil
	})
	return bucket
}

func (m BucketModel) All() (buckets Buckets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bucketId []byte, b *bolt.Bucket) error {
			m := getMeta(b)
			buckets = append(buckets, Bucket{Id: bucketId, Meta: m})
			return nil
		})
		return nil
	})
	return buckets, err
}
