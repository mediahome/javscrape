package scrape

import (
	"errors"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/query"
)

const javbusCNURL = "https://www.javbus.com/%s"
const javbusJAURL = "https://www.javbus.com/ja/%s"
const javbusENURL = "https://www.javbus.com/en/%s"
const javbusKOURL = "https://www.javbus.com/ko/%s"
const javbusUncensored = "javbusUncensored/search/%s&type=1"
const javbusCensored = "search/%s&type=1&parent=uc"

var grabJavbusLanguageList = []string{
	LanguageChinese:  javbusCNURL,
	LanguageEnglish:  javbusENURL,
	LanguageJapanese: javbusJAURL,
	LanguageKorea:    javbusKOURL,
}

type grabJAVBUS struct {
	language GrabLanguage
	doc      *goquery.Document
}

// Find ...
func (g *grabJAVBUS) Find(name string) (IGrab, error) {
	ug := *g
	url := grabJavbusLanguageList[g.language]
	results, e := g.getIndex(url, name)
	if e != nil {
		return nil, e
	}
	log.Println(results)
	return &ug, nil
}

type javbusSearchResult struct {
	Uncensored  bool
	Title       string
	PhotoFrame  string
	PhotoInfo   string
	ID          string
	ReleaseDate string
}

func (g *grabJAVBUS) getIndex(url string, name string) ([]*javbusSearchResult, error) {
	document, e := query.New(fmt.Sprintf(fmt.Sprintf(url, javbusCensored), name))
	isUncensored := false
	if e != nil {
		document, e = query.New(fmt.Sprintf(fmt.Sprintf(url, javbusUncensored), name))
		if e != nil {
			return nil, e
		}
		isUncensored = true
	}
	return javbusSearchResultAnalyze(document, isUncensored)
}

func javbusSearchResultAnalyze(document *goquery.Document, b bool) ([]*javbusSearchResult, error) {
	var res []*javbusSearchResult
	document.Find("#waterfall > a.movie-box").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Html())
	})
	if res == nil || len(res) == 0 {
		return nil, errors.New("no data found")
	}
	return res, nil
}

// Decode ...
func (g *grabJAVBUS) Decode(msg *Message) error {
	msg.Title = g.doc.Find("div.container > h3").Text()
	movie := g.doc.Find("div.container > div.row.movie")
	movie.Find("div.info > p").Each(func(i int, selection *goquery.Selection) {
		selection.Find("span").Each(func(i int, selection *goquery.Selection) {
			log.Println("index", i, "text", selection.Text())
		})
	})
	return nil
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
