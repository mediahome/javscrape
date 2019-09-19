package scrape

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/query"
)

const javbusCNURL = "https://www.javbus.com/%s"
const javbusJAURL = "https://www.javbus.com/ja/%s"
const javbusENURL = "https://www.javbus.com/en/%s"
const javbusKOURL = "https://www.javbus.com/ko/%s"
const uncensored = "uncensored/%s"

var grabJavbusLanguageList = []string{
	LanguageChinese:  javbusCNURL,
	LanguageEnglish:  javbusENURL,
	LanguageJapanese: javbusJAURL,
	LanguageKorea:    javbusKOURL,
}

type grabJAVBUS struct {
	language     GrabLanguage
	isUncensored bool
	doc          *goquery.Document
}

// Find ...
func (g *grabJAVBUS) Find(name string) (IGrab, error) {
	ug := *g
	url := grabJavbusLanguageList[g.language]
	document, e := query.New(fmt.Sprintf(url, name))
	if e != nil {
		document, e = query.New(fmt.Sprintf(fmt.Sprintf(url, uncensored), name))
		if e != nil {
			return nil, e
		}
		ug.isUncensored = true
	}
	ug.doc = document
	//ret, e := document.Html()
	//log.Println(ret)
	return &ug, nil
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
