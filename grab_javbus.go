package scrape

import (
	"errors"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/query"
)

const javbusCNURL = "https://www.javbus.com/"
const javbusJAURL = "https://www.javbus.com/ja/"
const javbusENURL = "https://www.javbus.com/en/"
const javbusKOURL = "https://www.javbus.com/ko/"
const javbusUncensored = "uncensored/search/%s&type=1"
const javbusCensored = "search/%s&type=1"

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
	for _, r := range results {
		log.Printf("%+v", r)
	}
	return &ug, nil
}

type javbusSearchResult struct {
	Uncensored  bool
	DetailLink  string
	Title       string
	PhotoFrame  string
	PhotoInfo   string
	ID          string
	ReleaseDate string
}

func (g *grabJAVBUS) getIndex(url string, name string) ([]*javbusSearchResult, error) {
	searchURL := fmt.Sprintf(url+javbusCensored, name)
	document, e := query.New(searchURL)
	isUncensored := false
	if e != nil {
		searchURL = fmt.Sprintf(url+javbusUncensored, name)
		document, e = query.New(searchURL)
		if e != nil {
			return nil, e
		}
		isUncensored = true
	}
	return javbusSearchResultAnalyze(document, isUncensored)
}

func javbusSearchResultAnalyze(document *goquery.Document, b bool) ([]*javbusSearchResult, error) {
	var res []*javbusSearchResult

	document.Find("#waterfall > div > a.movie-box").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javbusSearchResult)
		resTmp.Uncensored = b
		link, b := selection.Attr("href")
		if b {
			log.Println("link:", link)
			resTmp.DetailLink = link
		}
		src, b := selection.Find("#waterfall > div > a.movie-box > div.photo-frame > img").Attr("src")
		if b {
			resTmp.PhotoFrame = src
		}
		title, b := selection.Find("#waterfall > div > a.movie-box > div.photo-frame > img").Attr("title")
		if b {
			resTmp.Title = title
		}
		selection.Find("#waterfall > div > a.movie-box > div.photo-info > span > date").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				resTmp.ID = selection.Text()
			} else if i == 1 {
				resTmp.ReleaseDate = selection.Text()
			} else {
				//todo
			}

		})
		res = append(res, resTmp)
	})
	if res == nil || len(res) == 0 {
		return nil, errors.New("no data found")
	}
	return res, nil
}

// Decode ...
func (g *grabJAVBUS) Decode(msg *Message) error {

	return nil
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
