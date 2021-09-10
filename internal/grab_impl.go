package internal

import (
	"errors"

	"github.com/goextension/gomap"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/internal/action"
	"github.com/javscrape/go-scrape/rule"
)

var ErrActionIsAlreadyExist = errors.New("action is already exist")

type grabImpl struct {
	core.IScrape
	mainPage  string
	entrance  string
	inputType rule.InputType
	inputKey  string
	actions   map[string]*action.Action
	group     map[string][]*action.Action
	value     gomap.Map
}

func (g *grabImpl) InputType() rule.InputType {
	return g.inputType
}

func (g *grabImpl) InputKey() string {
	return g.inputKey
}

func (g *grabImpl) Put(key string, value *core.Value) {
	g.value.Set(key, value)
}

func (g *grabImpl) Get(key string) *core.Value {
	v := g.value.Get(key)
	return (v).(*core.Value)
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

func (g *grabImpl) Run(input string) error {
	return g.actionDo(g.entrance, input)
}

func (g *grabImpl) actionDo(name string, input string) error {
	actions := g.getActions(name)
	if len(actions) == 0 {
		return nil
	}
	log.Debug("GRAB", "start action", name, "query", input)
	for _, a := range actions {
		if err := a.Run(input); err != nil {
			return g.actionDo(a.Failure(), input)
		}
		return g.actionDo(a.Success(), input)
	}
	return nil
}

func (g *grabImpl) getActions(name string) []*action.Action {
	var exist bool
	var actions action.Actions
	if _, exist = g.actions[name]; exist {
		actions = []*action.Action{g.actions[name]}
	} else if _, exist = g.group[name]; exist {
		actions = g.group[name]
	}
	return actions.Sort()
}

func (g *grabImpl) Value() gomap.Map {
	return g.value
}

func NewGrab(scrape core.IScrape, r *rule.Rule) core.IGrab {
	return &grabImpl{
		IScrape:   scrape,
		mainPage:  r.MainPage,
		entrance:  r.Entrance,
		inputType: r.InputType,
		inputKey:  r.InputKey,
		actions:   make(map[string]*action.Action),
		group:     make(map[string][]*action.Action),
		value:     gomap.New(),
	}
}

var _ core.IGrab = (*grabImpl)(nil)
