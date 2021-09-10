package core

import (
	"github.com/javscrape/go-scrape/cache"
)

// IGrab ...
type IGrab interface {
	MainPage() string
	Cache() cache.Querier
	Do() error
}
