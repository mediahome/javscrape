package action

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
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
	log.Debug("ACTION", "query selector", a.action.Web.Selector)
	if a.action.Web.Selector != "" {
		find := query.Find(a.action.Web.Selector)
		return find.Html()
	}
	return query.Html()
}
