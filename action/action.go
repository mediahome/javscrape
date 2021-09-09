package action

import (
	"github.com/javscrape/go-scrape/rule"
)

type Action struct {
	name      string
	index     int
	filter    string
	uri       string
	through   bool
	onSuccess string
	onFailure string
}

func (a Action) Name() string {
	return a.name
}

func (a Action) Index() int {
	return a.index
}

func (a Action) Filter() string {
	return a.filter
}

func (a Action) Uri() string {
	return a.uri
}

func (a Action) Through() bool {
	return a.through
}

func (a Action) OnSuccess() string {
	return a.onSuccess
}

func (a Action) OnFailure() string {
	return a.onFailure
}

func (a Action) Do() {

}

func FromAction(action rule.Action) *Action {
	return &Action{
		name:      action.Name,
		index:     action.Index,
		filter:    action.Filter,
		uri:       action.URI,
		through:   action.Through,
		onSuccess: action.OnSuccess,
		onFailure: action.OnFailure,
	}
}
