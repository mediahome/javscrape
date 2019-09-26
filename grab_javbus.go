package scrape

import (
	"errors"
	"fmt"
	"strings"

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
	res      []*javbusSearchResult
}

// Name ...
func (g *grabJAVBUS) Name() string {
	return "javbus"
}

// Decode ...
func (g *grabJAVBUS) Decode([]*Message) error {
	for _, r := range g.res {
		detail, e := javbusSearchDetailAnalyze(r)
		if e != nil {
			return e
		}
		log.Infof("javbus detail:%+v", detail)
	}
	return nil
}

// Find ...
func (g *grabJAVBUS) Find(name string) (IGrab, error) {
	name = strings.ToUpper(name)
	ug := *g
	url := grabJavbusLanguageList[g.language]
	results, e := javbusSearchResultAnalyze(url, name)
	if e != nil {
		return nil, e
	}
	if debug {
		for _, r := range results {
			log.Infof("%+v", r)
		}
	}
	ug.res = results
	return &ug, nil
}

type javbusSearchResult struct {
	Uncensored bool
	DetailLink string
	Title      string
	PhotoFrame string
	//PhotoInfo   string
	ID          string
	ReleaseDate string
}

func javbusSearchResultAnalyze(url, name string) ([]*javbusSearchResult, error) {
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

	var res []*javbusSearchResult
	document.Find("#waterfall > div > a.movie-box").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javbusSearchResult)
		resTmp.Uncensored = isUncensored
		link, b := selection.Attr("href")
		if b {
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

type javbusSearchDetail struct {
}

func javbusSearchDetailAnalyze(result *javbusSearchResult) (*javbusSearchDetail, error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javbus search result is null")
	}
	document, e := query.New(result.DetailLink)
	if e != nil {
		return nil, e
	}

	title := document.Find("body > div.container > h3").Text()
	log.With("title", title).Info(result.ID)
	bigImage, exists := document.Find("body > div.container > div.row.movie > div > a.bigImage").Attr("href")
	log.With("bigImage", bigImage).Info(exists)
	image, exists := document.Find("body > div.container > div.row.movie > div > a > img").Attr("src")
	log.With("image", image).Info(exists)
	bigTitle, exists := document.Find("body > div.container > div.row.movie > div > a > img").Attr("title")
	log.With("bigTitle", bigTitle).Info(exists)
	document.Find("body > div.container > div.row.movie > div.col-md-3.info > p").Each(func(i int, selection *goquery.Selection) {
		switch i {
		case 0:
			id := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			log.With("id", id).Info("movie")
		case 1:
			date := goquery.NewDocumentFromNode(selection.Contents().Nodes[1]).Text()
			log.With("release date", date).Info("movie")
		case 2:
			length := goquery.NewDocumentFromNode(selection.Contents().Nodes[1]).Text()
			log.With("length", length).Info("movie")
		case 3:
			director := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			log.With("director", director).Info("movie")
		case 4:
			studio := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			log.With("studio", studio).Info("movie")
		case 5:
			label := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			log.With("label", label).Info("movie")
		case 6:
			series := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			log.With("series", series).Info("movie")
		case 7:
			//genre := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
			//log.With("genre", genre).Info("movie")
		case 8:
			selection.Find("span.genre > a").Each(func(i int, selection *goquery.Selection) {
				log.With("genre", selection.Text()).Info("movie")
			})
		case 9:
		case 10:
			selection.Find("span.genre > a").Each(func(i int, selection *goquery.Selection) {
				val, _ := selection.Attr("href")
				log.With("idols", selection.Text(), "image", val).Info("movie")
			})
			//.col-md-3 > p:nth-child(9) > span:nth-child(1) > a:nth-child(1)
		}
		log.With("index", i, "text", selection.Text()).Info("movie info")
		selection.Contents().Each(func(i int, selection *goquery.Selection) {
			log.With("content", selection.Text()).Info("contents")
		})
	})

	return &javbusSearchDetail{}, nil
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
