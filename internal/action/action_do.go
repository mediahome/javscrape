package action

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

func (a Action) Run() error {
	web, err := a.doWeb()
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

func (a Action) getWebURL(relative bool) string {
	value := a.getWebValue()
	mainPage := a.Get("main_page").GetString()
	if relative {
		if mainPage == "" {
			mainPage = value
		} else {
			mainPage = core.URL(mainPage, value)
		}
	} else {
		mainPage = core.URL(value)
	}
	return mainPage
}

func (a Action) getWebValue() string {
	var ret string
	if len(a.action.Web.Value) != 0 {
		var exps []string
		var vals []interface{}
		for _, s := range a.action.Web.Value {
			val := s[1:]
			switch s[0] {
			case '$':
				vals = append(vals, a.Get(val).GetString())
			case '%':
				exps = append(exps, val)
			default:
				vals = append(vals, val)
			}
		}
		if len(exps) == 0 {
			exps = append(exps, "%v")
		}
		format := strings.Join(exps, "/")
		fix := strings.Count(format, "%") - len(vals)
		for ; fix > 0; fix-- {
			vals = append(vals, "")
		}
		ret = fmt.Sprintf(format, vals...)
		log.Debug("ACTION", "get from value", ret)
	}
	return ret
}

func (a Action) doWeb() (sl string, err error) {
	log.Debug("ACTION", "do web query")
	var query *goquery.Document

	//url = core.URL(url, a.action.Web.BeforeURL)
	//mainSkipped := isSkipped(rule.SkipTypeMainPage, a.action.Web.Skip)
	//if mainSkipped {
	//	url = ""
	//}

	url := a.getWebURL(a.action.Web.Relative)
	if !isSkipped(rule.SkipTypeInput, a.action.Web.Skip) {
		url = a.getInputURL(url, a.Get(a.InputKey()).GetString())

	}
	log.Debug("ACTION", "query page url", url)
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
			if v != nil {
				log.Debug("ACTION", "put web value", "name", s.Name, "value", v, "index", i)
				a.Put(s.Name, v)
			}
		}
	}
}

func (a *Action) doWebSuccessValue(selection *goquery.Selection, p rule.Process) *core.Value {
	var ret *core.Value
	switch p.Property {
	case "attr":
		v := selection.AttrOr(p.PropertyName, "")
		if p.Trim {
			v = strings.TrimSpace(v)
		}
		ret = core.NewStringValue(v)
	}

	return ret
}

func (a Action) doWebSuccessPut() {

}
