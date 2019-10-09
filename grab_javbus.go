package scrape

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/net"
)

// DefaultJavbusMainPage ...
const DefaultJavbusMainPage = "https://www.javbus.com"
const javbusCNURL = "/"
const javbusJAURL = "/ja/"
const javbusENURL = "/en/"
const javbusKOURL = "/ko/"
const javbusUncensored = "uncensored/search/%s"
const javbusCensored = "search/%s"

var grabJavbusLanguageList = []string{
	LanguageChinese:  javbusCNURL,
	LanguageEnglish:  javbusENURL,
	LanguageJapanese: javbusJAURL,
	LanguageKorea:    javbusKOURL,
}

type grabJavbus struct {
	mainPage string
	sample   bool
	language GrabLanguage
	details  []*javbusSearchDetail
}

// Clone ...
func (g *grabJavbus) Clone() IGrab {
	panic("implement me")
}

// HasNext ...
func (g *grabJavbus) HasNext() bool {
	panic("implement me")
}

// Next ...
func (g *grabJavbus) Next() (IGrab, error) {
	panic("implement me")
}

// MainPage ...
func (g *grabJavbus) MainPage(url string) {
	g.mainPage = url
}

// sample ...
func (g *grabJavbus) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabJavbus) Name() string {
	return "javbus"
}

// Decode ...
func (g *grabJavbus) Decode(msg *[]*Content) error {
	for idx, detail := range g.details {
		log.With("index", idx).Info("decode")
		*msg = append(*msg, &Content{
			ID:            detail.id,
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
	return nil
}

// Find ...
func (g *grabJavbus) Find(name string) (IGrab, error) {
	ug := *g

	name = strings.ToUpper(name)
	url := ug.mainPage + grabJavbusLanguageList[ug.language]
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
		detail, e := javbusSearchDetailAnalyze(&ug, r)
		if e != nil {
			log.Error(e)
			continue
		}
		detail.uncensored = r.Uncensored
		detail.thumbImage = r.PhotoFrame
		detail.title = r.Title
		ug.details = append(ug.details, detail)
		log.Infof("javbus detail:%+v", detail)
	}

	return &ug, nil
}

type javbusSearchResult struct {
	Uncensored  bool
	DetailLink  string
	Title       string
	PhotoFrame  string
	ID          string
	ReleaseDate string
}

func javbusSearchResultAnalyze(url, name string) ([]*javbusSearchResult, error) {
	searchURL := fmt.Sprintf(url+javbusCensored, name)
	document, e := net.Query(searchURL)
	isUncensored := false
	if e != nil {
		searchURL = fmt.Sprintf(url+javbusUncensored, name)
		document, e = net.Query(searchURL)
		if e != nil {
			return nil, e
		}
		isUncensored = true
	}

	var res []*javbusSearchResult
	document.Find("#waterfall > div > a.movie-box").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javbusSearchResult)
		resTmp.Uncensored = isUncensored
		resTmp.DetailLink, _ = selection.Attr("href")
		resTmp.PhotoFrame, _ = selection.Find("#waterfall > div > a.movie-box > div.photo-frame > img").Attr("src")
		resTmp.Title, _ = selection.Find("#waterfall > div > a.movie-box > div.photo-frame > img").Attr("title")
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
	var idols []*Star
	if debug {
		log.Info(selection.Next().Html())
	}

	selection.Next().Find("div.star-box.idol-box").Each(func(i int, selection *goquery.Selection) {
		starLink := selection.Find("li > a").AttrOr("href", "")
		image := selection.Find("li > a > img").AttrOr("src", "")
		name := selection.Find("li > div.star-name > a").Text()
		name = strings.TrimSpace(name)
		log.With("name", name, "image", image, "star", starLink).Info("idols")
		idols = append(idols, &Star{
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
	var genre []*Genre
	selection.Next().Find("p > span.genre > a").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text()).Info("genre")
		g := new(Genre)
		g.Content = strings.TrimSpace(selection.Text())
		g.URL = selection.AttrOr("href", "")
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

const javbusTimeFormat = "2006-01-02"

func javbusSearchDetailAnalyzeDate(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 1 {
		return errors.New("wrong date node size")
	}
	date := goquery.NewDocumentFromNode(nodes[1]).Text()
	if debug {
		log.With("release date", date).Info("movie")
	}
	parse, e := time.Parse(javbusTimeFormat, date)
	if e != nil {
		return e
	}
	detail.date = parse
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
func javbusSearchDetailAnalyze(grab *grabJavbus, result *javbusSearchResult) (*javbusSearchDetail, error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javbus search result is null")
	}
	document, e := net.Query(result.DetailLink)
	if e != nil {
		return nil, e
	}

	detail := &javbusSearchDetail{}
	//detail.title = document.Find("body > div.container > h3").Text()
	//log.With("title", detail.title).Info(result.ID)
	detail.bigImage = document.Find("body > div.container > div.row.movie > div > a > img").AttrOr("src", "")
	log.With("image", detail.bigImage).Info("movie")
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
			image := selection.AttrOr("href", "")
			thumb := selection.Find("div > img").AttrOr("src", "")
			title := selection.Find("div > img").AttrOr("title", "")
			if debug {
				log.With("index", i, "image", image, "title", title, "thumb", thumb).Info("sample")
			}
			detail.sample = append(detail.sample, &Sample{
				Index: i,
				Thumb: thumb,
				Image: image,
				Title: title,
			})
		})
	}

	return detail, nil
}

// GrabJavbusOptions ...
type GrabJavbusOptions func(javbus *grabJavbus)

// JavbusLang ...
func JavbusLang(language GrabLanguage) GrabJavbusOptions {
	return func(javbus *grabJavbus) {
		javbus.language = language
	}
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(ops ...GrabJavbusOptions) IGrab {
	grab := &grabJavbus{
		mainPage: DefaultJavbusMainPage,
		language: LanguageJapanese,
	}
	for _, op := range ops {
		op(grab)
	}
	return grab
}
