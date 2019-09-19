package scrape

import (
	"fmt"
	"log"

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
	language GrabLanguage
}

// Find ...
func (g *grabJAVBUS) Find(name string) error {
	url := grabJavbusLanguageList[g.language]
	document, e := query.New(fmt.Sprintf(url, name))
	if e != nil {
		document, e = query.New(fmt.Sprintf(fmt.Sprintf(url, uncensored), name))
		if e != nil {
			return e
		}
	}
	ret, e := document.Html()
	if e != nil {
		return e
	}
	log.Println(ret)
	return nil
}

// Decode ...
func (g *grabJAVBUS) Decode(*Message) error {
	panic("implement me")
}

// NewGrabJAVBUS ...
func NewGrabJAVBUS(language GrabLanguage) IGrab {
	return &grabJAVBUS{
		language: language,
	}
}
