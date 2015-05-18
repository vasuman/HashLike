package db

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/boltdb/bolt"
)

var (
	logger          *log.Logger
	db              *bolt.DB
	bucKeyGroup     = []byte("Group")
	bucKeyChallenge = []byte("Challenge")
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

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

func encJSON(v interface{}) error {
	
	return nil
}

func setupBuckets(tx *bolt.Tx) error {
	var err error
	createBucket := func(bucket []byte) {
		if err == nil {
			_, err = tx.CreateBucketIfNotExists(bucket)
		}
	}
	createBucket(bucKeyGroup)
	createBucket(bucKeyChallenge)
	return err
}

func Init(dbInst *bolt.DB, logInst *log.Logger) error {
	db = dbInst
	logger = logInst
	err := db.Update(setupBuckets)
	return err
}
