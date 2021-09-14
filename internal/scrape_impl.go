package internal

import (
	"fmt"

	"github.com/goextension/log"
	"github.com/google/uuid"

	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/config"
	"github.com/javscrape/go-scrape/core"
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
	fmt.Println("DEBUG ON:", s.config.Debug)
	core.InitGlobalLogger(s.config.Debug)
}

func (s *scrapeImpl) Cache() cache.Querier {
	return s.cache
}

func (s *scrapeImpl) LoadRules(rs ...*rule.Rule) ([]core.IGrab, error) {
	if len(rs) == 0 {
		return nil, core.ErrEmptyRule
	}
	var gs []core.IGrab
	for i := range rs {
		log.Debug("SCRAPE", "new grab", "index", i)
		if rs[i].InputType == "" {
			rs[i].InputType = rule.InputTypeURL
		}
		if rs[i].Name == "" {
			rs[i].Name = uuid.New().String()
		}
		g := NewGrab(s, rs[i])
		if err := g.LoadActions(rs[i].Actions...); err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}

var _ core.IScrape = (*scrapeImpl)(nil)
