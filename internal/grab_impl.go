package internal

import (
	"errors"
	"sync"

	"github.com/goextension/gomap"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/internal/action"
	"github.com/javscrape/go-scrape/rule"
)

var ErrActionIsAlreadyExist = errors.New("action is already exist")

type grabImpl struct {
	core.IScrape
	inputType rule.InputType
	inputKey  string
	actions   map[string]*action.Action
	group     map[string][]*action.Action
	value     struct {
		lock sync.RWMutex
		keys map[string]struct{}
		gomap.Map
	}
}

func (g *grabImpl) InputType() rule.InputType {
	return g.inputType
}

func (g *grabImpl) InputKey() string {
	return g.inputKey
}

func (g *grabImpl) Name() string {
	return g.Get("#name").GetString()
}

func (g *grabImpl) PutInner(key string, value *core.Value) {
	g.value.lock.Lock()
	g.value.Set("#"+key, value)
	g.value.lock.Unlock()
}

func (g *grabImpl) Put(key string, value *core.Value) {
	g.value.lock.Lock()
	g.value.keys[key] = core.Empty
	g.value.Set(key, value)
	g.value.lock.Unlock()
}

func (g *grabImpl) Get(key string) *core.Value {
	var v interface{}
	g.value.lock.RLock()
	v = g.value.Get(key)
	g.value.lock.RUnlock()
	if v == nil {
		return nil
	}
	return (v).(*core.Value)
}

func (g *grabImpl) Value() gomap.Map {
	g.value.lock.Lock()
	defer g.value.lock.Unlock()
	var keys []string
	for k := range g.value.keys {
		keys = append(keys, k)
	}
	return g.value.Only(keys)
}

func (g *grabImpl) MainPage() string {
	return g.Get("#main_page").GetString()
}

func (g *grabImpl) Entrance() string {
	return g.Get("#entrance").GetString()
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
	g.PutInner(g.InputKey(), core.NewStringValue(input))
	return g.actionDo(g.Entrance())
}

func (g *grabImpl) actionDo(name string) error {
	actions := g.getActions(name)
	log.Debug("GRAB", "get actions", name, "total", len(actions))
	if len(actions) == 0 {
		return nil
	}
	log.Debug("GRAB", "start action", name, "query", g.Get(g.InputKey()))
	for _, a := range actions {
		if err := a.Run(); err != nil {
			return g.actionDo(a.Failure())
		}
		return g.actionDo(a.Success())
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

func NewGrab(scrape core.IScrape, r *rule.Rule) core.IGrab {
	value := gomap.New()
	for s, i := range r.Preset {
		value.Set(s, i)
	}
	value.Set("#name", core.NewStringValue(r.Name))
	if r.MainPage != "" {
		value.Set("#main_page", core.NewStringValue(r.MainPage))
	}
	if r.Entrance != "" {
		value.Set("#entrance", core.NewStringValue(r.Entrance))
	}

	if r.InputKey == "" {
		r.InputKey = "intput"
	}

	value.Range(func(key string, value interface{}) bool {
		v := value.(*core.Value)
		log.Debug("GRAB", "init map value", key, v)
		return true
	})

	return &grabImpl{
		IScrape:   scrape,
		inputType: r.InputType,
		inputKey:  r.InputKey,
		actions:   make(map[string]*action.Action),
		group:     make(map[string][]*action.Action),
		value: struct {
			lock sync.RWMutex
			keys map[string]struct{}
			gomap.Map
		}{Map: value, keys: make(map[string]struct{})},
	}
}

var _ core.IGrab = (*grabImpl)(nil)
