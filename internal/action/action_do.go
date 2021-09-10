package action

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

func (a Action) Run(input string) error {
	web, err := a.doWeb(a.MainPage(), input)
	if err != nil {
		return err
	}
	log.Debug("ACTION", "web html", web)
	return nil
}

func (a Action) getInputURL(urlpath string, input string) string {
	switch a.InputType() {
	case rule.InputTypeURL:
		return core.URL(urlpath, input)
	case rule.InputTypeValue:
		return core.URLAddValues(urlpath, url.Values{
			a.InputKey(): []string{input},
		})
	}
	return ""
}

func isSkipped(skipType rule.SkipType, skips []rule.SkipType) bool {
	if len(skips) == 0 {
		return false
	}
	for _, skip := range skips {
		if skip == skipType {
			return true
		}
	}
	return false
}

func (a Action) doWeb(url string, input string) (sl string, err error) {
	var query *goquery.Document

	//url = core.URL(url, a.action.Web.BeforeURL)
	mainSkipped := isSkipped(rule.SkipTypeMainPage, a.action.Web.Skip)
	if mainSkipped {
		url = ""
	}

	if a.action.Web.BeforeURL != "" {
		if a.action.Web.Relative {
			if url == "" {
				url = a.action.Web.BeforeURL
			} else {
				url = core.URL(url, a.action.Web.BeforeURL)
			}
		} else {
			if url != "" {
				return "", core.ErrAbsoluteMultiAddress
			} else {
				url = a.action.Web.BeforeURL
			}
		}
	}

	log.Debug("ACTION", "get main url", url)
	if len(a.action.Web.FromValue) != 0 {
		var froms []string
		for _, s := range a.action.Web.FromValue {
			froms = append(froms, a.Get(s).GetString())
		}
		if len(froms) == 1 {
			if a.action.Web.Relative {
				url = core.URL(url, froms...)
			} else {
				url = core.URL(froms[0])
			}
		} else if len(froms) > 1 {
			url = core.URL(url, froms...)
			if !a.action.Web.Relative {
				return "", errors.New("absolute mode cannot use multi from")
			}
		} else {
			//0
		}
		log.Debug("ACTION", "get from value", url)
	}

	if a.action.Web.AfterURL != "" {
		url = core.URL(url, a.action.Web.AfterURL)
		log.Debug("ACTION", "get page after url", url)
	}

	if !isSkipped(rule.SkipTypeInput, a.action.Web.Skip) {
		url = a.getInputURL(url, input)
		log.Debug("ACTION", "query page url", url)
	}

	query, err = a.Cache().Query(url, false)

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
