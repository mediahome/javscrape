package action

import (
	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

type Action struct {
	core.IGrab
	action *rule.Action
}

func FromAction(grab core.IGrab, action rule.Action) *Action {
	return &Action{
		IGrab:  grab,
		action: &action,
	}
}

func (a *Action) Failure() string {
	return a.action.OnFailure
}

func (a *Action) Success() string {
	return a.action.OnSuccess
}

func (a Action) GetValue(key string) (core.KeyType, string) {
	val := key[1:]
	switch key[0] {
	case '$':
		return core.KeyTypeCache, a.Get(val).GetString()
	case '%':
		return core.KeyTypeExpression, key
	case '#':
		return core.KeyTypeSystem, a.Get(key).GetString()
	}
	return core.KeyTypeProto, key
}
