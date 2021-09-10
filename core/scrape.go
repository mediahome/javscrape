package core

import (
	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/rule"
)

// IScrape ...
type IScrape interface {
	Cache() cache.Querier
	LoadRules(r ...*rule.Rule) ([]IGrab, error)
}
