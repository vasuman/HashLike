package cache

import (
	"sync"

	"github.com/vasuman/HashLike/models"
)

var (
	cLock sync.RWMutex
	cMap  map[string]float64
)

func GetCount(url, sys string) (float64, error) {
	cLock.RLock()
	v, ok := cMap[url+":-:"+sys]
	cLock.RUnlock()
	if ok {
		return v
	}
	l, err := models.GetLocation(url)
	if err == nil {

	} else if err == models.ErrNoSuchLocation {

	}
}

func Invalidate(url, sys string) {
	cLock.Lock()
	delete(cMap, url+":-:"+sys)
	cLock.Unlock()
}
