package db

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vasuman/HashLike/pow"
)

var (
	ErrNoSuchChallenge = errors.New("no such challenge")
)

type Challenge struct {
	Challenge []byte
	ChParams  *pow.ChallengeParams
	Expires   time.Time
}

func SaveChallenge(c *Challenge) error {
	err := db.Update(func(tx *bolt.Tx) error {
		challengeBucket := tx.Bucket(bucKeyChallenge)
		v, err := encGob(c)
		if err != nil {
			return err
		}
		return challengeBucket.Put([]byte(c.Challenge), v)
	})
	return err
}

func GetChallenge(challenge []byte) (*Challenge, error) {
	c := new(Challenge)
	err := db.View(func(tx *bolt.Tx) error {
		challengeBucket := tx.Bucket(bucKeyChallenge)
		v := challengeBucket.Get(challenge)
		if v == nil {
			return ErrNoSuchChallenge
		}
		return decGob(v, c)
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func runFlush(tx *bolt.Tx) error {
	return nil
}

func flushChallenges(ticker *time.Ticker) {
	ok := true
	for ok {
		select {
		case _, ok = <-ticker.C:
			err := db.Update(runFlush)
			if err != nil {
				logger.Printf("error flushing %v", err)
			}
		}
	}
}
