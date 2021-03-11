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

// DefaultJavbusMainPage ...
const DefaultJavbusMainPage = "https://www.javbus.com"
const javbusCNURL = "/"
const javbusJAURL = "/ja/"
const javbusENURL = "/en/"
const javbusKOURL = "/ko/"
const javbusUncensored = "uncensored/search/%s"
const javbusCensored = "search/%s"

var grabJavbusLanguageList = []string{
	LanguageChineseTraditional: javbusCNURL,
	LanguageEnglish:            javbusENURL,
	LanguageJapanese:           javbusJAURL,
	LanguageKorea:              javbusKOURL,
}

type grabJavbus struct {
	mainPage   string
	next       string
	uncensored bool
	sample     bool
	exact      bool
	finder     string
	language   GrabLanguage
	details    []*javbusSearchDetail
	cache      *Cache
}

func (g *grabJavbus) SetLanguage(language GrabLanguage) {
	g.language = language
}

// SetExact ...
func (g *grabJavbus) SetExact(b bool) {
	g.exact = b
}

// SetSample ...
func (g *grabJavbus) SetSample(b bool) {
	g.sample = b
}

// HasNext ...
func (g *grabJavbus) HasNext() bool {
	return g.next != ""
}

// Next ...
func (g *grabJavbus) Next() (IGrab, error) {
	return g.find(g.next)
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
func (g *grabJavbus) Result() (c []Content, e error) {
	for idx, detail := range g.details {
		if debug {
			log.Infow("decode", "index", idx, "id", detail.id)
		}

		c = append(c, Content{
			From:          g.Name(),
			Language:      g.language.String(),
			Uncensored:    detail.uncensored,
			ID:            strings.ToUpper(detail.id),
			Title:         detail.title,
			OriginalTitle: "",
			Year:          strconv.Itoa(detail.date.Year()),
			Poster:        detail.bigImage,
			Thumb:         detail.thumbImage,
			ReleaseDate:   detail.date,
			Studio:        detail.studio,
			MovieSet:      detail.series,
			Director:      detail.director,
			Publisher:     detail.label,
			Plot:          "",
			Genres:        detail.genre,
			Actors:        detail.idols,
			Sample:        detail.sample,
		})
	}
	return
}
func (g *grabJavbus) clone() *grabJavbus {
	clone := new(grabJavbus)
	*clone = *g
	clone.details = nil
	return clone
}

func (g *grabJavbus) find(url string) (IGrab, error) {
	clone := g.clone()
	results, e := javbusSearchResultAnalyze(clone, url)
	if e != nil {
		return clone, e
	}
	for _, r := range results {
		if debug {
			log.Infow("find", "id", r.ID, "detail", r)
		}

		if clone.exact && strings.ToLower(r.ID) != strings.ToLower(clone.finder) {
			continue
		}
		detail, e := javbusSearchDetailAnalyze(clone, r)
		if e != nil {
			log.Error(e)
			continue
		}
		detail.uncensored = r.Uncensored
		detail.thumbImage = r.PhotoFrame
		detail.title = r.Title
		clone.details = append(clone.details, detail)
		if debug {
			log.Infow("find|detail", "id", detail.id, "detail", detail)
		}
	}

	return clone, nil
}

// Find ...
func (g *grabJavbus) Find(name string) (IGrab, error) {
	g.finder = name
	url := g.mainPage + grabJavbusLanguageList[g.language]
	g.uncensored = false
	grab, e := g.find(fmt.Sprintf(url+javbusCensored, name))
	if e != nil {
		g.uncensored = true
		return g.find(fmt.Sprintf(url+javbusUncensored, name))
	}
	return grab, nil
}

type javbusSearchResult struct {
	Uncensored  bool
	DetailLink  string
	Title       string
	PhotoFrame  string
	ID          string
	ReleaseDate string
}

func javbusSearchResultAnalyze(grab *grabJavbus, url string) ([]*javbusSearchResult, error) {
	document, e := grab.cache.Query(url)
	if e != nil {
		return nil, e
	}

	var res []*javbusSearchResult
	document.Find("#waterfall > div > a.movie-box").Each(func(i int, selection *goquery.Selection) {
		resTmp := new(javbusSearchResult)
		resTmp.Uncensored = grab.uncensored
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
	next, b := document.Find("body > div.text-center.hidden-xs > ul > li > a#next").Attr("href")
	//if debug {
	log.Infow("pagination", "next", next, "exist", b)
	//}
	grab.next = ""
	if b && next != "" {
		grab.next = grab.mainPage + next
	}
	if res == nil || len(res) == 0 {
		return nil, errors.New("no data found")
	}
	return res, nil
}

type javbusSearchDetail struct {
	language   string
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
	//javbusSearchDetailAnalyzeDummy,
}

var analyzeLanguageList = map[GrabLanguage][]string{
	LanguageEnglish: {
		"ID",
		"Release Date",
		"Length",
		"Director",
		"Studio",
		"Label",
		"Series",
		"Genre",
		"JAV Idols",
	},
	LanguageJapanese: {
		"品番",
		"発売日",
		"収録時間",
		"監督",
		"メーカー",
		"レーベル",
		"シリーズ",
		"ジャンル",
		"出演者",
	},
	LanguageChineseTraditional: {
		"識別碼",
		"發行日期",
		"長度",
		"導演",
		"製作商",
		"發行商",
		"系列",
		"類別",
		"演員",
	},
}

func getAnalyzeLanguageFunc(language GrabLanguage, selection *goquery.Selection) AnalyzeLanguageFunc {
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	text = strings.TrimSpace(text)
	if text == "" {
		return javbusSearchDetailAnalyzeDummy
	}
	for idx, list := range analyzeLanguageList[language] {
		ret := strings.Index(text, list)
		if debug {
			log.Infow("LanguageFunc", "ret", ret, "text", text, "list", list)
		}
		if ret != -1 {
			return analyzeLangFuncList[idx]
		}
	}
	return javbusSearchDetailAnalyzeDummy
}
func javbusSearchDetailAnalyzeDummy(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	log.Warnw("dummy", "text", text, "detail", detail, "size", len(selection.Contents().Nodes))
	return nil
}

//tag:idols
func javbusSearchDetailAnalyzeIdols(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	var idols []*Star
	selection.Next().Find("div.star-box.idol-box").Each(func(i int, selection *goquery.Selection) {
		starLink := selection.Find("li > a").AttrOr("href", "")
		image := selection.Find("li > a > img").AttrOr("src", "")
		name := selection.Find("li > div.star-name > a").Text()
		name = strings.TrimSpace(name)
		if debug {
			log.Infow("idols", "name", name, "image", image, "star", starLink)
		}
		idols = append(idols, &Star{
			StarLink: starLink,
			Image:    image,
			Name:     name,
		})
	})
	if debug {
		log.Infow("movie", "idols", idols)
	}
	detail.idols = idols
	return
}

//tag:series
func javbusSearchDetailAnalyzeSeries(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	series := ""
	if len(nodes) <= 2 {
		return errors.New("wrong series node size")
	}
	series = goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.Infow("movie", "series", series)
	}
	detail.series = series
	return
}

//tag:genre
func javbusSearchDetailAnalyzeGenre(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	var genre []*Genre
	selection.Next().Find("p > span.genre > label > a").Each(func(i int, selection *goquery.Selection) {
		if debug {
			log.Infow("genre", "text", selection.Text())
		}
		g := new(Genre)
		g.Content = strings.TrimSpace(selection.Text())
		g.URL = selection.AttrOr("href", "")
		genre = append(genre, g)
	})
	if debug {
		log.Infow("movie", "genre", genre)
	}
	detail.genre = genre
	return
}

//tags:label
func javbusSearchDetailAnalyzeLabel(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	if len(nodes) <= 2 {
		return errors.New("wrong label node size")
	}
	label := goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.Infow("movie", "label", label)
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
		log.Infow("movie", "studio", studio)
	}
	detail.studio = strings.TrimSpace(studio)
	return
}
func javbusSearchDetailAnalyzeDirector(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	nodes := selection.Contents().Nodes
	director := ""
	if len(nodes) <= 2 {
		return errors.New("wrong director node size")
	}
	director = goquery.NewDocumentFromNode(nodes[2]).Text()
	if debug {
		log.Infow("movie", "director", director)
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
		log.Infow("movie", "length", length)
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
		log.Infow("movie", "release date", date)
	}
	parse, e := time.Parse(javbusTimeFormat, strings.TrimSpace(date))
	if e != nil {
		return e
	}
	detail.date = parse
	return
}
func javbusSearchDetailAnalyzeID(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	if debug {
		html, e := selection.Html()
		log.Infow("AnalyzeID", "source", html, "error", e)
	}
	nodes := selection.Find("span").Contents().Nodes
	if len(nodes) < 2 {
		return errors.New("wrong id node size")
	}
	id := goquery.NewDocumentFromNode(nodes[1]).Text()
	if debug {
		log.Infow("movie", "id", id)
	}
	detail.id = strings.TrimSpace(id)
	return
}
func javbusSearchDetailAnalyze(grab *grabJavbus, result *javbusSearchResult) (*javbusSearchDetail, error) {
	if result == nil || result.DetailLink == "" {
		return nil, errors.New("javbus search result is null")
	}
	document, e := grab.cache.Query(result.DetailLink)
	if e != nil {
		return nil, e
	}

	detail := &javbusSearchDetail{}
	//detail.title = document.Find("body > div.container > h3").Text()
	//log.With("title", detail.title).Info(result.ID)
	detail.bigImage = document.Find("body > div.container > div.row.movie > div > a > img").AttrOr("src", "")
	if debug {
		log.Infow("movie", "image", detail.bigImage)
	}
	//detail.bigImage, exists = document.Find("body > div.container > div.row.movie > div > a.bigImage").Attr("href")
	//log.With("bigImage", detail.bigImage).Info(exists)
	//detail.title, exists = document.Find("body > div.container > div.row.movie > div > a > img").Attr("title")
	//log.With("bigTitle", detail.title).Info(exists)

	document.Find("body > div.container > div.row.movie > div.col-md-3.info > p").Each(func(i int, selection *goquery.Selection) {
		if debug {
			html, e := selection.Html()
			log.Infow("AnalyzeLanguageFunc", "language", grab.language, "source", html, "error", e)
			selection.Contents().Each(func(i int, selection *goquery.Selection) {
				log.Infow("AnalyzeLanguageFunc|Contents", "content", selection.Text())
			})
		}
		err := getAnalyzeLanguageFunc(grab.language, selection)(selection, detail)
		if err != nil {
			log.Error(err)
		}
	})

	if grab.sample {
		document.Find("#sample-waterfall > a.sample-box").Each(func(i int, selection *goquery.Selection) {
			image := selection.AttrOr("href", "")
			thumb := selection.Find("div > img").AttrOr("src", "")
			title := selection.Find("div > img").AttrOr("title", "")
			if debug {
				log.Infow("sample", "index", i, "image", image, "title", title, "thumb", thumb)
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

// JavbusExact ...
func JavbusExact(b bool) GrabJavbusOptions {
	return func(javbus *grabJavbus) {
		javbus.exact = b
	}
}

// NewGrabJavbus ...
func NewGrabJavbus(ops ...GrabJavbusOptions) IGrab {
	grab := &grabJavbus{
		cache:    NewCache(),
		mainPage: DefaultJavbusMainPage,
		language: LanguageChineseTraditional,
		exact:    true,
	}
	for _, op := range ops {
		op(grab)
	}
	return grab
}
