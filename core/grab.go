package core

import (
	"github.com/goextension/gomap"

	"github.com/javscrape/go-scrape/cache"
	"github.com/javscrape/go-scrape/rule"
)

// IGrab ...
type IGrab interface {
	MainPage() string
	LoadActions(...rule.Action) error
	Cache() cache.Querier
	InputType() rule.InputType
	InputKey() string
	Put(key string, value *Value)
	Get(key string) *Value
	Run(input string) error
	Value() gomap.Map
}

var Empty = struct{}{}
