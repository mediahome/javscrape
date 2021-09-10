package action

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

func (a Action) Do(key string) error {
	web, err := a.doWeb(a.MainPage(), key)
	if err != nil {
		return err
	}
	log.Debug("ACTION", "web html", web)
	return nil
}

func (a Action) doWeb(url string, key string) (sl string, err error) {
	var query *goquery.Document
	if a.action.Web.URI != "" {
		url = core.URL(url, a.action.Web.URI, key)
		log.Debug("ACTION", "query page uri", url)
		query, err = a.Cache().Query(url, false)
	}

	if a.action.Web.URL != "" {
		url = core.URL(a.action.Web.URL, key)
		log.Debug("ACTION", "query page url", url)
		query, err = a.Cache().Query(url, false)
	}

	if err != nil {
		return "", err
	}

	if query == nil {
		return "", nil
	}
	if a.action.Web.Selector != "" {
		log.Debug("ACTION", "do query selector", a.action.Web.Selector)
		find := query.Find(a.action.Web.Selector)
		a.doWebSuccess(find)
		return find.Html()
	}
	return query.Html()
}

func (a *Action) doWebSuccess(selection *goquery.Selection) {
	for i, s := range a.action.Web.Success {
		switch s.Type {
		case rule.ProcessTypePut:
			v := a.doWebSuccessValue(selection, s)
			log.Debug("ACTION", "put web value", "name", s.Name, "value", v, "index", i)
			a.Put(s.Name, v)
		}
	}
}

func (a *Action) doWebSuccessValue(selection *goquery.Selection, p rule.Process) *core.Value {
	var ret core.Value
	switch p.Property {
	case "attr":
		ret.Type = rule.ProcessValueString
		v := selection.AttrOr(p.PropertyName, "")
		if p.Trim {
			v = strings.TrimSpace(v)
		}
		ret.Set(v)
	}

	return &ret
}

func (a Action) doWebSuccessPut() {

}
