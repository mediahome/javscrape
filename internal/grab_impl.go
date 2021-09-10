package internal

import (
	"errors"

	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/internal/action"
	"github.com/javscrape/go-scrape/rule"
)

var ErrActionIsAlreadyExist = errors.New("action is already exist")

type grabImpl struct {
	core.IScrape
	mainPage string
	entrance string
	actions  map[string]*action.Action
	group    map[string][]*action.Action
}

func NewGrab(scrape core.IScrape, r *rule.Rule) core.IGrab {
	return &grabImpl{
		IScrape:  scrape,
		mainPage: r.MainPage,
		entrance: r.Entrance,
		actions:  make(map[string]*action.Action),
		group:    make(map[string][]*action.Action),
	}
}

func (g *grabImpl) MainPage() string {
	return g.mainPage
}

func (g *grabImpl) LoadActions(acts ...rule.Action) error {
	for _, v := range acts {
		switch v.Type {
		case rule.ActionTypeGroup:
			log.Debug("GRAB", "load group", v.Name)
			g.group[v.Name] = append(g.group[v.Name], action.FromAction(g, v))
		default:
			v.Type = rule.ActionTypeAction
			fallthrough
		case rule.ActionTypeAction:
			log.Debug("GRAB", "load action", v.Name)
			if _, exist := g.actions[v.Name]; exist {
				return ErrActionIsAlreadyExist
			}
			g.actions[v.Name] = action.FromAction(g, v)
		}
	}
	return nil
}

func (g *grabImpl) Do() error {
	actions := g.getEntranceActions()
	for _, a := range actions {
		if err := a.Do(); err != nil {
			return err
		}
	}
	return nil
}

func (g *grabImpl) getEntranceActions() []*action.Action {
	var exist bool
	var actions action.Actions
	if _, exist = g.actions[g.entrance]; exist {
		actions = []*action.Action{g.actions[g.entrance]}
	} else if _, exist = g.group[g.entrance]; exist {
		actions = g.group[g.entrance]
	}
	return actions.Sort()
}

var _ core.IGrab = (*grabImpl)(nil)
