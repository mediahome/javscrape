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
		ssel := selection.Clone()

		if s.Selector != "" {
			ssel = ssel.Find(s.Selector)
		}
		//html, _ := ssel.Html()
		//log.Debug("ACTION", "print find html", html)
		for _, f := range s.Filter {
			ssel = ssel.Filter(f)
		}

		if len(ssel.Nodes) == 0 || len(ssel.Nodes) < s.Index {
			log.Error("failed do loop by index", "loop", i, "index", s.Index, "length", len(ssel.Nodes), "name", s.Name)
			continue
		}
		//log.Debug("ACTION", "print filter html", html)
		ssel = goquery.NewDocumentFromNode(ssel.Nodes[s.Index]).First()
		html, _ := ssel.Html()
		log.Debug("ACTION", "print current html", "index", s.Index, html)
		//ssel.Each(func(i int, selection *goquery.Selection) {
		//	if s.Index == i {
		//		html, _ = ssel.Html()
		//		log.Debug("ACTION", "print each html", "index", i, html)
		//		ssel = selection
		//	}
		//})

		switch s.Type {
		case rule.ProcessTypePut:
			v := a.doWebSuccessValue(ssel, s)
			if v != nil {
				log.Debug("ACTION", "put web value", "name", s.Name, "value", v, "index", i)
				a.Put(s.Name, v)
			}
		}
	}
}

func (a *Action) doWebSuccessValue(selection *goquery.Selection, p rule.Process) *core.Value {
	var v string
	switch p.Property {
	case "array":
		var arr []interface{}
		selection.Each(func(i int, selection *goquery.Selection) {
			v = strings.TrimSpace(selection.Text())
			arr = append(arr, v)
		})
		if len(arr) != 0 {
			return core.NewArrayValue(arr)
		}
	case "value":
		v = selection.Text()
	case "attr":
		v = selection.AttrOr(p.PropertyName, "")
	case "text":
		selection.Contents().Each(func(i int, selection *goquery.Selection) {
			if goquery.NodeName(selection) == "#text" {
				v = selection.Text()
			}
		})
	}

	if p.Trim {
		v = strings.TrimSpace(v)
	}

	if v == "" {
		return nil
	}
	return core.NewStringValue(v)
}

func (a Action) doWebSuccessPut() {

}
