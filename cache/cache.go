package main

import (
	"sync"
	"time"
)

type ccItem struct {
	challenge string
	expiry    time.Time
	solved    bool
}

type cCache struct {
	expList []ccItem

	l sync.Mutex
	m map[string]ccItem
}

func (cc *cCache) flushOld() {
	cc.l.Lock()
	now := time.Now()
	for i, v := range cc.expList {
		if !now.After(v.expiry) {
		}
	}
	cc.l.Unlock()
}

func (cc *cCache) Get(challenge string) {

}
