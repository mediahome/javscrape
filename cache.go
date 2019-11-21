package scrape

import (
	_ "github.com/gocacher/badger-cache/easy"
)

type cacheImpl struct {
}

var cache *cacheImpl

func init() {
	cache = newCache()
}

func newCache() *cacheImpl {
	return &cacheImpl{}
}
