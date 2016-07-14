package models

import (
	"encoding/json"
	"time"

	"../lib"

	"github.com/boltdb/bolt"
)

type MetaBucket struct {
	Name      string `json:"name"`
	TimeStamp int    `json:"timestamp"`
	// TODO: more metadata
}

type Bucket struct {
	Id   string     `json:"id"`
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
		bucket.Id = string(bucketId)
		bucket.Meta = getMeta(b)
		return nil
	})
	return bucket
}

func (m BucketModel) All() (buckets Buckets, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bucketId []byte, b *bolt.Bucket) error {
			m := getMeta(b)
			buckets = append(buckets, Bucket{Id: string(bucketId), Meta: m})
			return nil
		})
		return nil
	})
	return buckets, err
}

func (m BucketModel) Create(bucketName string) (Bucket, error) {

	meta := MetaBucket{Name: bucketName, TimeStamp: int(time.Now().UTC().UnixNano() / 1000000)}
	bucket := Bucket{Id: lib.RandSeq(8), Meta: meta}

	err := db.Update(func(tx *bolt.Tx) error {
		b, cerr := tx.CreateBucket([]byte(bucket.Id))
		if cerr != nil {
			return cerr
		}

		j, jerr := json.Marshal(bucket.Meta)
		if jerr != nil {
			return jerr
		}

		perr := b.Put([]byte("meta"), j)
		if perr != nil {
			return perr
		}
		return nil
	})
	return bucket, err
}

func (m BucketModel) Destroy(bucketId string) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		derr := tx.DeleteBucket([]byte(bucketId))
		return derr
	})
	return err
}

func (m BucketModel) Set(bucket Bucket) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.Id))

		j, jerr := json.Marshal(bucket.Meta)
		if jerr != nil {
			return jerr
		}

		e := b.Put([]byte("meta"), j)
		if e != nil {
			return e
		}
		return e
	})
	return err
}
