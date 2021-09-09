package scrape

import (
	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/network"
)

// Options ...
type Options func(impl *scrapeImpl)

// CacheOption ...
func CacheOption(cache cache.Querier) Options {
	return func(impl *scrapeImpl) {
		impl.cache = cache
	}
}

func ProxyOption(addr string) Options {
	return func(impl *scrapeImpl) {
		err := network.RegisterProxy(addr)
		if err != nil {
			panic(err)
		}
	}
}
