package scrape

import (
	"github.com/goextension/log/zap"

	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/network"
	"github.com/javscrape/go-scrape/rule"
)

// IScrape ...
type IScrape interface {
	Cache() cache.Querier
	LoadAction(r rule.Rule) IGrab
}

// NewScrape ...
func NewScrape(opts ...Options) IScrape {
	scrape := &scrapeImpl{}

	for _, opt := range opts {
		opt(scrape)
	}

	scrape.init()
	return scrape
}

func init() {
	zap.InitZapSugar()
}

func (impl *scrapeImpl) init() {
	if impl.cache == nil {
		impl.cache = cache.NewQueryCache(network.Client())
	}
}
