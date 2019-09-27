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
	sample   bool
	language GrabLanguage
	details  []*javbusSearchDetail
}

// Sample ...
func (g *grabJAVBUS) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabJAVBUS) Name() string {
	return "javbus"
}

// Decode ...
func (g *grabJAVBUS) Decode([]*Message) error {

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
	for _, r := range results {
		detail, e := javbusSearchDetailAnalyze(g, r)
		if e != nil {
			continue
		}
		g.details = append(g.details, detail)
		log.Infof("javbus detail:%+v", detail)
	}

	if g.sample {

	}

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
	title      string
	thumbImage string
	bigImage   string
	id         string
	date       string
	length     string
	director   string
	studio     string
	label      string
	series     string
	genre      []string
	idols      []*Idols
	Sample     []*Sample
}

// AnalyzeLanguageFunc ...
type AnalyzeLanguageFunc func(selection *goquery.Selection, detail *javbusSearchDetail) (e error)

var analyzeLangFuncList = []AnalyzeLanguageFunc{
	javbusSearchDetailAnalyzeID,
	javbusSearchDetailAnalyzeDate,
	javbusSearchDetailAnalyzeLength,
	javbusSearchDetailAnalyzeDirector,
	javbusSearchDetailAnalyzeStudio,
	javbusSearchDetailAnalyzeLabel,
	javbusSearchDetailAnalyzeSeries,
	javbusSearchDetailAnalyzeGenre,
	javbusSearchDetailAnalyzeIdols,
	javbusSearchDetailAnalyzeDummy,
	javbusSearchDetailAnalyzeDummy,
}

var analyzeLanguageList = map[GrabLanguage][]string{
	LanguageEnglish: {
		"ID:",
		"Release Date:",
		"Length:",
		"Director:",
		"Studio:",
		"Label:",
		"Series:",
		"Genre:",
		"JAV Idols",
	},
	LanguageJapanese: {
		"品番:",
		"発売日:",
		"収録時間:",
		"監督:",
		"メーカー:",
		"レーベル:",
		"ジャンル:",
		"出演者",
	},
}

func getAnalyzeLanguageFunc(language GrabLanguage, selection *goquery.Selection) AnalyzeLanguageFunc {
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	for idx, list := range analyzeLanguageList[language] {
		if strings.Compare(text, list) == 0 {
			return analyzeLangFuncList[idx]
		}
	}
	return javbusSearchDetailAnalyzeDummy
}
func javbusSearchDetailAnalyzeDummy(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	log.With("size", len(selection.Contents().Nodes), "text", text).Warnf("%+v", *detail)
	return nil
}
func javbusSearchDetailAnalyzeIdols(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	var idols []*Idols
	log.Info(selection.Next().Html())

	selection.Next().Find("div.star-box.idol-box").Each(func(i int, selection *goquery.Selection) {
		starLink, _ := selection.Find("li > a").Attr("href")
		image, _ := selection.Find("li > a > img").Attr("src")
		name := selection.Find("li > div.star-name > a").Text()
		name = strings.TrimSpace(name)
		log.With("name", name, "image", image, "star", starLink).Info("idols")
		idols = append(idols, &Idols{
			StarLink: starLink,
			Image:    image,
			Name:     name,
		})
	})
	if debug {
		log.With("idols", idols).Info("movie")
	}
	detail.idols = idols
	return
}
func javbusSearchDetailAnalyzeSeries(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong director node size")
	}
	series := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.With("series", series).Info("movie")
	}
	detail.series = series
	return
}
func javbusSearchDetailAnalyzeGenre(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	var genre []string
	selection.Next().Find("p > span.genre > a").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text()).Info("genre")
		g := selection.Text()
		g = strings.TrimSpace(g)
		genre = append(genre, g)
	})
	if debug {
		log.With("genre", genre).Info("movie")
	}
	detail.genre = genre
	return
}
func javbusSearchDetailAnalyzeLabel(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong label node size")
	}
	label := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.With("label", label).Info("movie")
	}
	detail.label = strings.TrimSpace(label)
	return
}
func javbusSearchDetailAnalyzeStudio(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong studio node size")
	}
	studio := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.With("studio", studio).Info("movie")
	}
	detail.studio = strings.TrimSpace(studio)
	return
}
func javbusSearchDetailAnalyzeDirector(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong director node size")
	}
	director := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.With("director", director).Info("movie")
	}
	detail.director = strings.TrimSpace(director)
	return
}
func javbusSearchDetailAnalyzeLength(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 1 {
		return errors.New("wrong length node size")
	}
	length := goquery.NewDocumentFromNode(nodes[1]).Text()
	if debug {
		log.With("length", length).Info("movie")
	}
	detail.length = strings.TrimSpace(length)
	return
}
func javbusSearchDetailAnalyzeDate(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 1 {
		return errors.New("wrong date node size")
	}
	date := goquery.NewDocumentFromNode(nodes[1]).Text()
	if debug {
		log.With("release date", date).Info("movie")
	}
	detail.date = strings.TrimSpace(date)
	return
}
func javbusSearchDetailAnalyzeID(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong id node size")
	}
	id := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.With("id", id).Info("movie")
	}
	detail.id = strings.TrimSpace(id)
	return
}
func javbusSearchDetailAnalyze(grab *grabJAVBUS, result *javbusSearchResult) (*javbusSearchDetail, error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javbus search result is null")
	}
	document, e := query.New(result.DetailLink)
	if e != nil {
		return nil, e
	}

	detail := &javbusSearchDetail{}
	var exists bool
	detail.title = document.Find("body > div.container > h3").Text()
	log.With("title", detail.title).Info(result.ID)
	detail.bigImage, exists = document.Find("body > div.container > div.row.movie > div > a > img").Attr("src")
	log.With("image", detail.bigImage).Info(exists)
	//detail.bigImage, exists = document.Find("body > div.container > div.row.movie > div > a.bigImage").Attr("href")
	//log.With("bigImage", detail.bigImage).Info(exists)
	//detail.title, exists = document.Find("body > div.container > div.row.movie > div > a > img").Attr("title")
	//log.With("bigTitle", detail.title).Info(exists)

	document.Find("body > div.container > div.row.movie > div.col-md-3.info > p").Each(func(i int, selection *goquery.Selection) {
		err := getAnalyzeLanguageFunc(grab.language, selection)(selection, detail)
		if err != nil {
			log.Error(err)
		}
		if debug {
			log.With("index", i, "text", selection.Text()).Info("info movie")
			selection.Contents().Each(func(i int, selection *goquery.Selection) {
				log.With("content", selection.Text()).Info("info contents")
			})
		}
	})

	if grab.sample {
		document.Find("#sample-waterfall > a.sample-box").Each(func(i int, selection *goquery.Selection) {
			image, _ := selection.Attr("href")
			thumb, _ := selection.Find("div > img").Attr("src")
			title, _ := selection.Find("div > img").Attr("title")
			if debug {
				log.With("index", i, "image", image, "title", title, "thumb", thumb).Info("sample")
			}
			detail.Sample = append(detail.Sample, &Sample{
				Index: i,
				Thumb: thumb,
				Image: image,
				Title: title,
			})
		})
	}

	return detail, nil
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
