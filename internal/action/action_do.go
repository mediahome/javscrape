package action

import (
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
)

func (a Action) Do() error {
	web, err := a.doWeb(a.MainPage())
	if err != nil {
		return nil
	}
	if core.DEBUG {
		html, err := web.Html()
		if err != nil {
			return err
		}
		log.Debug("WEB", html)
	}
	return nil
}

func (a Action) doWeb(url string) (sl *goquery.Selection, err error) {
	var query *goquery.Document
	if a.action.Web.URI != "" {
		url = path.Join(url, a.action.Web.URI)
		query, err = a.Cache().Query(url, false)
	}

	if a.action.Web.URL != "" {
		query, err = a.Cache().Query(url, false)
	}

	if err != nil {
		return nil, err
	}

	return query.Find(a.action.Web.Selector).Find(a.action.Web.Selector), nil
}
