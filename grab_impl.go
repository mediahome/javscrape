package scrape

import (
	"errors"

	"github.com/javscrape/go-scrape/action"
	"github.com/javscrape/go-scrape/rule"
)

var ErrActionIsAlreadyExist = errors.New("action is already exist")

type grabImpl struct {
	mainPage string
	entrance string
	actions  map[string]*action.Action
	group    map[string][]*action.Action
}

func (g *grabImpl) MainPage() string {
	return g.mainPage
}

func (g *grabImpl) LoadAction(acts ...rule.Action) error {
	for _, v := range acts {
		switch v.Type {
		case rule.ActionTypeAction:
			if _, exist := g.actions[v.Name]; exist {
				return ErrActionIsAlreadyExist
			}
			g.actions[v.Name] = action.FromAction(v)
		case rule.ActionTypeActionGroup:
			g.group[v.Name] = append(g.group[v.Name], action.FromAction(v))
		}
	}
	return nil
}
