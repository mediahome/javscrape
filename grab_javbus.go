package scrape

import (
	"fmt"
	"log"

	"github.com/javscrape/go-scrape/query"
)

const javbusCN_URL = "https://www.javbus.com/%s"
const javbusJA_URL = "https://www.javbus.com/ja/%s"
const javbusEN_URL = "https://www.javbus.com/en/%s"

var grabJavbusLanguageList = []string{
	LanguageChinese:  javbusCN_URL,
	LanguageEnglish:  javbusEN_URL,
	LanguageJapanese: javbusJA_URL,
}

type grabJAVBUS struct {
	language GrabLanguage
}

// Find ...
func (g *grabJAVBUS) Find(name string) error {
	url := grabJavbusLanguageList[g.language]
	document, e := query.New(fmt.Sprintf(url, name))
	if e != nil {
		return e
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
