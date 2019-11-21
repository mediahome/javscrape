package scrape

import (
	"crypto/sha256"
	"fmt"

	"github.com/gocacher/badger-cache"
	"github.com/gocacher/cacher"
)

type cacheImpl struct {
	cache cacher.Cacher
}

var _cache *cacheImpl

func init() {
	_cache = newCache()
}

func newCache() *cacheImpl {
	return &cacheImpl{
		cache: cache.NewBadgerCache("tmp"),
	}

}

// Hash ...
func Hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", sum256)
}
