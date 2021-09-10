package internal

import (
	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/config"
	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/log"
	"github.com/javscrape/go-scrape/network"
	"github.com/javscrape/go-scrape/rule"
)

type scrapeImpl struct {
	config *config.Config
	cache  cache.Querier
	//grabs []Grab
}

// NewScrape ...
func NewScrape(opts ...Options) core.IScrape {
	scrape := &scrapeImpl{}

	for _, opt := range opts {
		opt(scrape)
	}

	scrape.init()
	return scrape
}

func (s *scrapeImpl) init() {
	if s.cache == nil {
		s.cache = cache.NewQueryCache(network.Client())
	}
	if s.config == nil {
		s.config = config.DefaultConfig()
	}

	log.InitGlobalLogger(s.config.Debug)
}

func (s *scrapeImpl) Cache() cache.Querier {
	return s.cache
}

func (s *scrapeImpl) LoadRules(rs ...*rule.Rule) ([]core.IGrab, error) {
	if len(rs) == 1 {
		return nil, core.ErrEmptyRule
	}
	var gs []core.IGrab
	for i, _ := range rs {
		gs = append(gs, NewGrab(s, rs[i]))
	}
	return gs, nil
}

var _ core.IScrape = (*scrapeImpl)(nil)
