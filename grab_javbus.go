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
		detail, e := javbusSearchDetailAnalyze(g.language, r)
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

type Idols struct {
	Image string
	Name  string
}

type javbusSearchDetail struct {
	id       string
	date     string
	length   string
	director string
	studio   string
	label    string
	series   string
	genre    []string
	idols    []*Idols
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
	//text := selection.Find("p > span.header").Text()
	//genre := selection.Text()
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	for idx, list := range analyzeLanguageList[language] {
		if strings.Compare(text, list) == 0 {
			return analyzeLangFuncList[idx]
		}
		if debug {
			//log.With("source", list, "target", text).Info(strings.Compare(text, list))
		}
	}
	//6 is genre
	//if strings.Compare(analyzeLanguageList[language][6], genre) == 0 {
	//	return analyzeLangFuncList[6]
	//}
	//if debug {
	//	log.With("source", analyzeLanguageList[language][6], "target", genre).Info(strings.Compare(analyzeLanguageList[language][6], genre))
	//}
	return javbusSearchDetailAnalyzeDummy
}
func javbusSearchDetailAnalyzeDummy(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	text := goquery.NewDocumentFromNode(selection.Contents().Nodes[0]).Text()
	log.With("size", len(selection.Contents().Nodes), "text", text).Warnf("%+v", *detail)
	return nil
}
func javbusSearchDetailAnalyzeIdols(selection *goquery.Selection, detail *javbusSearchDetail) (e error) {
	var idols []*Idols
	selection.Next().Find("span.genre > a").Each(func(i int, selection *goquery.Selection) {
		val, _ := selection.Attr("href")
		name := selection.Text()
		name = strings.TrimSpace(name)
		log.With("text", selection.Text(), "url", val).Info("idols")
		idols = append(idols, &Idols{
			Image: val,
			Name:  name,
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
func javbusSearchDetailAnalyze(lan GrabLanguage, result *javbusSearchResult) (*javbusSearchDetail, error) {
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
	deatil := &javbusSearchDetail{}
	document.Find("body > div.container > div.row.movie > div.col-md-3.info > p").Each(func(i int, selection *goquery.Selection) {
		err := getAnalyzeLanguageFunc(lan, selection)(selection, deatil)
		if err != nil {
			log.Error(err)
		}
		//switch i {
		//case 0:
		//	id := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	log.With("id", id).Info("movie")
		//case 1:
		//	date := goquery.NewDocumentFromNode(selection.Contents().Nodes[1]).Text()
		//	log.With("release date", date).Info("movie")
		//case 2:
		//	length := goquery.NewDocumentFromNode(selection.Contents().Nodes[1]).Text()
		//	log.With("length", length).Info("movie")
		//case 3:
		//	director := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	log.With("director", director).Info("movie")
		//case 4:
		//	studio := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	log.With("studio", studio).Info("movie")
		//case 5:
		//	label := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	log.With("label", label).Info("movie")
		//case 6:
		//	series := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	log.With("series", series).Info("movie")
		//case 7:
		//	//genre := goquery.NewDocumentFromNode(selection.Contents().Nodes[2]).Text()
		//	//log.With("genre", genre).Info("movie")
		//case 8:
		//	selection.Find("span.genre > a").Each(func(i int, selection *goquery.Selection) {
		//		log.With("genre", selection.Text()).Info("movie")
		//	})
		//case 9:
		//case 10:
		//	selection.Find("span.genre > a").Each(func(i int, selection *goquery.Selection) {
		//		val, _ := selection.Attr("href")
		//		log.With("idols", selection.Text(), "image", val).Info("movie")
		//	})
		//	//.col-md-3 > p:nth-child(9) > span:nth-child(1) > a:nth-child(1)
		//}
		if debug {
			log.With("index", i, "text", selection.Text()).Info("info movie")
			selection.Contents().Each(func(i int, selection *goquery.Selection) {
				log.With("content", selection.Text()).Info("info contents")
			})
		}
	})

	return &javbusSearchDetail{}, nil
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
