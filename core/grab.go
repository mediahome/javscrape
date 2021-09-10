package core

import (
	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/rule"
)

// IGrab ...
type IGrab interface {
	MainPage() string
	LoadActions(...rule.Action) error
	Cache() cache.Querier
	Put(key, value string)
	Get(key string) string
	Do(key string) error
}
