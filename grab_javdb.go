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
const DefaultJavdbMainPage = "https://javdb4.com"
const javdbSearch = "/search?q=%s&f=all"

const javdbENRUL = "locale=en"
const javdbZHRUL = "locale=zh"

const javdbZHNo = "番號"
const javdbZHTime = "日期"
const javdbZHTimeLong = "時長"
const javdbZHDirector = "導演"
const javdbZHStudio = "片商"
const javdbZHPublisher = "發行"
const javdbZHIdols = "演員"
const javdbZHGenre = "類別"
const javdbZHSeries = "系列"
const javdbZHTimeFormat = "2006-01-02"
const javdbZHRating = "評分"

type grabJavdb struct {
	mainPage string
	next     string
	sample   bool
	exact    bool
	finder   string
	details  []*javdbSearchDetail
	cache    *Cache
}

// SetExact ...
func (g *grabJavdb) SetExact(b bool) {
	g.exact = b
}

// SetSample ...
func (g *grabJavdb) SetSample(b bool) {
	g.sample = b
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

// SetSample ...
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
			log.Infof("javdb detail:%+v", detail)
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
	rating     string
}

func javdbSearchDetailAnalyze(grab *grabJavdb, result *javdbSearchResult) (detail *javdbSearchDetail, e error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javdb search result is null")
	}
	document, e := grab.cache.Query(grab.mainPage + result.DetailLink)
	if e != nil {
		return nil, e
	}
	detail = new(javdbSearchDetail)
	detail.uncensored = strings.Index(document.Find("h2.title > strong").Text(), "無碼") > 0
	log.Infow("DetailAnalyze", "uncensored", detail.uncensored)

	detail.bigImage, _ = document.Find("div.columns > div.column > a > img").Attr("src")
	log.Infow("DetailAnalyze", "bigImage", detail.bigImage)
	document.Find("div.columns > div > nav.panel > div.panel-block").Each(func(i int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find("strong").Text())
		value := strings.TrimSpace(selection.Find("span.value").Text())
		if debug {
			log.Infow("DetailAnalyze", "title", title, "value", value)
		}
		switch {
		case strings.Index(title, javdbZHNo) != -1:
			detail.id = value
		case strings.Index(title, javdbZHTime) != -1:
			detail.date, _ = time.Parse(javdbZHTimeFormat, value)
		case strings.Index(title, javdbZHTimeLong) != -1:
			detail.length = value
		case strings.Index(title, javdbZHDirector) != -1:
			detail.director = value
		case strings.Index(title, javdbZHStudio) != -1:
			detail.studio = value
		case strings.Index(title, javdbZHSeries) != -1:
			detail.series = value
		case strings.Index(title, javdbZHRating) != -1:
			detail.rating = value
		case strings.Index(title, javdbZHPublisher) != -1:
			//nothing
		case strings.Index(title, javdbZHIdols) != -1:
			var idols []*Star
			selection.Find("span.value>a").Each(func(i int, selection *goquery.Selection) {
				s := &Star{}
				s.Name = strings.TrimSpace(selection.Text())
				s.StarLink = grab.mainPage + selection.AttrOr("href", "")
				idols = append(idols, s)
			})
			detail.idols = idols
		case strings.Index(title, javdbZHGenre) != -1:
			var genre []*Genre
			selection.Find("span.value>a").Each(func(i int, selection *goquery.Selection) {
				g := &Genre{}
				g.Content = strings.TrimSpace(selection.Text())
				g.URL = grab.mainPage + selection.AttrOr("href", "")
				genre = append(genre, g)
				if debug {
					log.Infow("DetailAnalyze|Genre", "context", g.Content, "url", g.URL)
				}
			})
			detail.genre = genre
		default:
			log.Warnw("DetailAnalyze|Other", "title", title, "value", value)
		}
	})

	if grab.sample {
		document.Find("div.message-body > div.tile-images.preview-images > a.tile-item").Each(func(i int, selection *goquery.Selection) {
			image := selection.AttrOr("href", "")
			title := selection.AttrOr("data-caption", "")
			detail.sample = append(detail.sample, &Sample{
				Index: i,
				Image: image,
				Title: title,
			})
			if debug {
				log.Infow("DetailAnalyze|Sample", "index", i, "image", image, "title", title, "thumb", "")
			}
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

//tag:SearchResultAnalyze
func javdbSearchResultAnalyze(grab *grabJavdb, url string) (result []*javdbSearchResult, e error) {
	document, e := grab.cache.Query(url)
	if e != nil {
		return nil, e
	}
	var res []*javdbSearchResult
	document.Find("#videos > div > div.grid-item.column").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javdbSearchResult)

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
		if debug {
			log.Infow("SearchResultAnalyze", "res", res)
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
func (g *grabJavdb) Result() (c []*Content, e error) {

	for idx, detail := range g.details {
		log.Infow("decode", "index", idx)
		c = append(c, &Content{
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
		})
	}
	return
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
		cache:    NewCache(),
	}
	for _, op := range ops {
		op(grab)
	}

	return grab
}
