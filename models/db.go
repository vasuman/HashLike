package models

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var (
	groupBucKey = []byte("Group")
)

func encGob(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func decGob(b []byte, dst interface{}) error {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	err := dec.Decode(dst)
	return err
}

func setupBuckets(tx *bolt.Tx) error {
	var err error
	createBucket := func(bucket []byte) {
		if err == nil {
			_, err = tx.CreateBucketIfNotExists(bucket)
		}
	}
	createBucket(groupBucKey)
	return err
}

func InitDb(dbInst *bolt.DB) error {
	db = dbInst
	err := db.Update(setupBuckets)
	return err
}
