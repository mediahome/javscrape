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
