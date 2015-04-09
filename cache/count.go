package cache

import (
	"sync"

	"github.com/vasuman/HashLike/models"
)

var (
	cntLock = new(sync.RWMutex)
	cntMap  = make(map[string]float64)
)

const countSep = ":=:"

func GetCount(url, sys string) (float64, error) {
	cntLock.RLock()
	v, ok := cntMap[url+countSep+sys]
	cntLock.RUnlock()
	if ok {
		return v, nil
	}
	l, err := models.GetLocation(url)
	if err != nil {
		return 0, err
	}
	// Lookup `solutions`
}

func InvalidateCount(url, sys string) {
	cntLock.Lock()
	delete(cntMap, url+":-:"+sys)
	cntLock.Unlock()
}
