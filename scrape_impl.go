package scrape

import (
	"github.com/javscrape/go-scrape/action"
	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/config"
	"github.com/javscrape/go-scrape/rule"
)

type scrapeImpl struct {
	config *config.Config
	cache  cache.Querier
}

func (impl *scrapeImpl) Cache() cache.Querier {
	return impl.cache
}

func (impl *scrapeImpl) LoadAction(r rule.Rule) (IGrab, error) {
	g := &grabImpl{
		mainPage: r.MainPage,
		entrance: r.Entrance,
		actions:  make(map[string]*action.Action),
		group:    make(map[string][]*action.Action),
	}
	if err := g.LoadAction(r.Actions...); err != nil {
		return nil, err
	}
	return g, nil
}
