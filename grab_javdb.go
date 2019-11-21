package scrape

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"
)

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const javdbSearch = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	next     string
	sample   bool
	exact    bool
	finder   string
	details  []*javdbSearchDetail
}

func (g *grabJavdb) clone() *grabJavdb {
	clone := new(grabJavdb)
	*clone = *g
	clone.details = nil
	//clone.details = make([]*javdbSearchDetail, len(g.details))
	//copy(clone.details, g.details)
	return clone
}

// HasNext ...
func (g *grabJavdb) HasNext() bool {
	return g.next != ""
}

// Next ...
func (g *grabJavdb) Next() (IGrab, error) {
	return g.find(g.next)
}

// Sample ...
func (g *grabJavdb) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabJavdb) Name() string {
	return "javdb"
}

func (g *grabJavdb) find(url string) (IGrab, error) {
	clone := g.clone()
	results, e := javdbSearchResultAnalyze(clone, url)
	if e != nil {
		return clone, e
	}
	if debug {
		for _, r := range results {
			log.Infof("%+v", r)
		}
	}

	for _, r := range results {
		if clone.exact && strings.ToLower(r.ID) != strings.ToLower(clone.finder) {
			log.Infow("continue", "id", r.ID, "find", clone.finder)
			continue
		}
		detail, e := javdbSearchDetailAnalyze(clone, r)
		if e != nil {
			log.Error(e)
			continue
		}
		detail.thumbImage = r.Thumb
		detail.title = r.Title
		clone.details = append(clone.details, detail)
		if debug {
			log.Infof("javbus detail:%+v", detail)
		}
	}

	return clone, nil
}

// Find ...
func (g *grabJavdb) Find(name string) (IGrab, error) {
	g.finder = name
	url := fmt.Sprintf(g.mainPage+javdbSearch, name)
	return g.find(url)
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

func javdbSearchDetailAnalyze(grab *grabJavdb, result *javdbSearchResult) (detail *javdbSearchDetail, e error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javdb search result is null")
	}
	document, e := Query(grab.mainPage + result.DetailLink)
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
				g.URL = grab.mainPage + selection.AttrOr("href", "")
				genre = append(genre, g)
			})
			detail.genre = genre
		case 5:
			var idols []*Star
			selection.Find("a").Each(func(i int, selection *goquery.Selection) {
				s := new(Star)
				s.Name = strings.TrimSpace(selection.Text())
				s.StarLink = grab.mainPage + selection.AttrOr("href", "")
				idols = append(idols, s)
			})
			detail.idols = idols
		}
	})

	if grab.sample {
		document.Find("div.message-body > div.tile-images.preview-images > a.tile-item").Each(func(i int, selection *goquery.Selection) {
			image := selection.AttrOr("href", "")
			//thumb := selection.Find("div > img").AttrOr("src", "")
			title := selection.AttrOr("data-caption", "")
			if debug {
				log.Infow("sample", "index", i, "image", image, "title", title, "thumb", "")
			}
			detail.sample = append(detail.sample, &Sample{
				Index: i,
				//Thumb: thumb,
				Image: image,
				Title: title,
			})
		})
	}

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

func javdbSearchResultAnalyze(grab *grabJavdb, url string) (result []*javdbSearchResult, e error) {
	document, e := Query(url)
	if e != nil {
		return nil, e
	}
	var res []*javdbSearchResult
	document.Find("#videos > div > div.grid-item.column").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javdbSearchResult)
		if debug {
			log.Infow("javdb", "index", i, "text", selection.Text())
		}
		//resTmp.Title, _ = selection.Find("a.box").Attr("Title")
		resTmp.DetailLink = selection.Find("a.box").AttrOr("href", "")

		resTmp.Thumb = selection.Find("a.box > div.item-image > img").AttrOr("src", "")
		if strings.Index(resTmp.Thumb, "//") == 0 {
			resTmp.Thumb = "https:" + resTmp.Thumb
		}
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

	next, b := document.Find("body > section > div > nav.pagination > a.pagination-next").Attr("href")
	if debug {
		log.Infow("pagination", "next", next, "exist", b)
	}
	grab.next = ""
	if b && next != "" {
		grab.next = grab.mainPage + next
	}

	if res == nil || len(res) == 0 {
		return nil, errors.New("no data found")
	}
	return res, nil
}

// Decode ...
func (g *grabJavdb) Decode(msg *Content) error {
	for idx, detail := range g.details {
		log.Infof("decode", "index", idx)
		*msg = Content{
			From:          g.Name(),
			Uncensored:    detail.uncensored,
			ID:            strings.ToUpper(detail.id),
			Title:         detail.title,
			OriginalTitle: "",
			Year:          strconv.Itoa(detail.date.Year()),
			Image:         detail.bigImage,
			Thumb:         detail.thumbImage,
			ReleaseDate:   detail.date,
			Studio:        detail.studio,
			MovieSet:      detail.series,
			Plot:          "",
			Genres:        detail.genre,
			Actors:        detail.idols,
			Sample:        detail.sample,
		}
	}
	return nil
}

// GrabJavdbOptions ...
type GrabJavdbOptions func(javdb *grabJavdb)

// JavdbExact ...
func JavdbExact(b bool) GrabJavdbOptions {
	return func(javdb *grabJavdb) {
		javdb.exact = b
	}
}

// MainPage ...
func (g *grabJavdb) MainPage(url string) {
	g.mainPage = url
}

// NewGrabJavdb ...
func NewGrabJavdb(ops ...GrabJavdbOptions) IGrab {
	grab := &grabJavdb{
		mainPage: DefaultJavdbMainPage,
		sample:   false,
		exact:    true,
	}
	for _, op := range ops {
		op(grab)
	}
	return grab
}
