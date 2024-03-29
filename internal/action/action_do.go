package action

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/gomap"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

func (a Action) Run() error {

	_, err := a.doWeb()
	if err != nil {
		return err
	}
	log.Debug("ACTION", "web html")
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
	mainPage := a.MainPage()
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
			t, v := a.GetValue(s)
			switch t {
			case core.KeyTypeExpression:
				exps = append(exps, v)
			default:
				vals = append(vals, v)
			}
		}
		format := "%v"
		if len(exps) == 1 {
			format = exps[0]
		} else if len(exps) > 1 {
			format = strings.Join(exps, "/")
		}

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
	webCache := gomap.New()

	url := a.getWebURL(a.action.Web.Relative)
	if !isSkipped(rule.SkipTypeInput, a.action.Web.Skip) {
		url = a.getInputURL(url, a.Get("#"+a.InputKey()).GetString())
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
		a.processDo(webCache, find, a.action.Web.Success)
		return find.Html()
	}
	return query.Html()
}

func (a *Action) processDo(cache gomap.Map, selection *goquery.Selection, process []rule.Process) *goquery.Selection {
	if len(process) == 0 {
		return selection
	}
	for i, s := range process {
		ssel := selection.Clone()
		if s.Selector != "" {
			ssel = ssel.Find(s.Selector)
		}
		ssel = a.processDo(cache, ssel, s.Compare)

		if ssel == nil {
			continue
		}
		html, _ := ssel.Html()
		log.Debug("ACTION", "print compare html", "index", s.Index, html)

		switch s.Type {
		case rule.ProcessTypeCompare:
			var ret *goquery.Selection
			selection.EachWithBreak(func(i int, selection *goquery.Selection) bool {
				v := core.ProcessValue(ssel, s)
				if v == nil {
					return true
				}
				if strings.Contains(v.GetString(), s.Name) {
					ret = selection
					return false
				}
				return true
			})
			return ret
		case rule.ProcessTypePutArray:
			v := core.ProcessValue(ssel, s)
			if v != nil {
				log.Debug("ACTION", "put web value", "name", s.Name, "value", v, "index", i)
				a.Put(s.Name, v)
			}
		case rule.ProcessTypePut:
			ssel = goquery.NewDocumentFromNode(ssel.Nodes[s.Index]).First()
			html, _ := ssel.Html()
			log.Debug("ACTION", "print current html", "index", s.Index, html)
			v := core.ProcessValue(ssel, s)
			if v != nil {
				log.Debug("ACTION", "put web value", "name", s.Name, "value", v, "index", i)
				a.Put(s.Name, v)
			}
		}
	}
	return nil
}

func (a Action) doWebSuccessPut() {

}
