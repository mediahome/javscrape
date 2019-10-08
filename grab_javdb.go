package scrape

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/net"
)

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const javdbSearch = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	sample   bool
	//details  []*javdbSearchDetail
}

// Sample ...
func (g *grabJavdb) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabJavdb) Name() string {
	return "javdb"
}

// Find ...
func (g *grabJavdb) Find(name string) (IGrab, error) {
	url := g.mainPage + javdbSearch
	results, e := javdbSearchResultAnalyze(url, name)
	if e != nil {
		return nil, e
	}
	if debug {
		for _, r := range results {
			log.Infof("%+v", r)
		}
	}
	return g, nil
}

type javdbSearchResult struct {
	ID         string
	Title      string
	DetailLink string
	Thumb      string
	Tags       []string
	Date       string
}

func javdbSearchResultAnalyze(url, name string) (result []*javdbSearchResult, e error) {
	document, e := net.Query(fmt.Sprintf(url, name))
	if e != nil {
		return nil, e
	}
	var res []*javdbSearchResult
	document.Find("#videos > div > div.grid-item.column").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javdbSearchResult)
		if debug {
			log.With("index", i, "text", selection.Text()).Info("javdb")
		}
		//resTmp.Title, _ = selection.Find("a.box").Attr("Title")
		resTmp.DetailLink, _ = selection.Find("a.box").Attr("href")
		resTmp.Thumb, _ = selection.Find("a.box > div.item-image > img").Attr("src")
		resTmp.ID = selection.Find("a.box > div.uid").Text()
		resTmp.Title = selection.Find("a.box >div.video-title").Text()
		selection.Find("a.box > div.tags > span.tag").Each(func(i int, selection *goquery.Selection) {
			resTmp.Tags = append(resTmp.Tags, selection.Text())
		})
		resTmp.Date = selection.Find("a.box >div.meta").Text()
		if resTmp.ID != "" {
			res = append(res, resTmp)
		}
	})
	if res == nil || len(res) == 0 {
		return nil, errors.New("no data found")
	}
	return res, nil
}

// Decode ...
func (g *grabJavdb) Decode(*[]*Message) error {
	panic("implement me")
}

// MainPage ...
func (g *grabJavdb) MainPage(url string) {
	g.mainPage = url
}

// NewJavdb ...
func NewJavdb() IGrab {
	return &grabJavdb{
		mainPage: DefaultJavdbMainPage,
		sample:   false,
	}
}
