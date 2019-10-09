package scrape

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/net"
)

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const javdbSearch = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	sample   bool
	details  []*javdbSearchDetail
}

// Clone ...
func (g *grabJavdb) Clone() IGrab {
	panic("implement me")
}

// HasNext ...
func (g *grabJavdb) HasNext() bool {
	panic("implement me")
}

// Next ...
func (g *grabJavdb) Next() (IGrab, error) {
	panic("implement me")
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
	ug := *g
	url := ug.mainPage + javdbSearch
	results, e := javdbSearchResultAnalyze(url, name)
	if e != nil {
		return nil, e
	}
	if debug {
		for _, r := range results {
			log.Infof("%+v", r)
		}
	}
	for _, r := range results {
		detail, e := javdbSearchDetailAnalyze(&ug, r)
		if e != nil {
			log.Error(e)
			continue
		}
		detail.thumbImage = r.Thumb
		detail.title = r.Title
		ug.details = append(ug.details, detail)
		log.Infof("javbus detail:%+v", detail)
	}

	return &ug, nil
}

type javdbSearchDetail struct {
	title      string
	thumbImage string
	bigImage   string
	id         string
	date       time.Time
	length     string
	director   string
	studio     string
	label      string
	series     string
	genre      []*Genre
	idols      []*Star
	sample     []*Sample
	uncensored bool
}

const javdbTimeFormat = "2006-01-02"

func javdbSearchDetailAnalyze(javdb *grabJavdb, result *javdbSearchResult) (detail *javdbSearchDetail, e error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javdb search result is null")
	}
	document, e := net.Query(javdb.mainPage + result.DetailLink)
	if e != nil {
		return nil, e
	}
	detail = new(javdbSearchDetail)
	detail.uncensored = strings.Index(document.Find("h2.title > strong").Text(), "無碼") > 0
	detail.bigImage, _ = document.Find("div.columns.item-content > div.column.column-video-cover > a > img").Attr("src")
	document.Find("div.columns.item-content > div > nav.panel > div.panel-block > span.value").Each(func(i int, selection *goquery.Selection) {
		switch i {
		case 0:
			detail.id = selection.Text()
		case 1:
			detail.date, _ = time.Parse(javdbTimeFormat, selection.Text())
		case 2:
			detail.length = selection.Text()
		case 3:
			detail.studio = selection.Text()
		case 4:
			var genre []*Genre
			selection.Find("a").Each(func(i int, selection *goquery.Selection) {
				g := new(Genre)
				g.Content = strings.TrimSpace(selection.Text())
				g.URL = javdb.mainPage + selection.AttrOr("href", "")
				genre = append(genre, g)
			})
			detail.genre = genre
		case 5:
			var idols []*Star
			selection.Find("a").Each(func(i int, selection *goquery.Selection) {
				s := new(Star)
				s.Name = strings.TrimSpace(selection.Text())
				s.StarLink = javdb.mainPage + selection.AttrOr("href", "")
				idols = append(idols, s)
			})
			detail.idols = idols
		}

	})
	return detail, nil
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
		resTmp.DetailLink = selection.Find("a.box").AttrOr("href", "")
		resTmp.Thumb = "https" + selection.Find("a.box > div.item-image > img").AttrOr("src", "")
		resTmp.ID = selection.Find("a.box > div.uid").Text()
		resTmp.Title = selection.Find("a.box >div.video-title").Text()
		selection.Find("a.box > div.tags > span.tag").Each(func(i int, selection *goquery.Selection) {
			resTmp.Tags = append(resTmp.Tags, selection.Text())
		})
		resTmp.Date = strings.TrimSpace(selection.Find("a.box >div.meta").Text())
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
func (g *grabJavdb) Decode(*[]*Content) error {
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
