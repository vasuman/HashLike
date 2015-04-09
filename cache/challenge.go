package main

import (
	"sync"
	"time"
)

type ccItem struct {
	url       string
	challenge []byte
	ip        []byte
	when      time.Time
	expiry    time.Time
	inv       bool
}

var (
	chLock  = new(sync.Mutex)
	expList = make([]*ccItem, 0)
	chMap   = make(map[string]*ccItem)
)

// Checks if challenge exists and is still valid. If so, then mark
// it as invalid
func CheckAndMark(challenge []byte) bool {
	k := string(challenge)
	chLock.Lock()
	defer chLock.Unlock()
	i, ok := chMap[k]
	if !ok || i.inv {
		return false
	}
	i.inv = true
	return true
}

// Assumes all that given expiry is always greater than that of all
// the pending challenges
func Add(challenge []byte, expiry time.Time) {
	k := string(challenge)
	c := &ccItem{
		challenge: challenge,
		expiry:    expiry,
		inv:       false,
	}
	chLock.Lock()
	chMap[k] = c
	append(expList, c)
	chLock.Unlock()
}

func flushExpired() {
	now := time.Now()
	i := 0
	chLock.Lock()
	for i, v := range expList {
		if !now.After(v.expiry) {
			break
		}
		delete(chMap, string(v.challenge))
	}
	expList = expList[i:]
	chLock.Unlock()
}
