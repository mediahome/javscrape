package scrape

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/javscrape/go-scrape/net"
)

// DefaulBp4xMainPage ...
const DefaulBp4xMainPage = "https://www.bp4x.com"
const bp4xJavURL = "/?q=%s"
const bp4xAmateurURL = "/?c=amateur&q=%s"
const bp4xIVURL = "/?c=iv&q=%s"
const bp4xHentaiURL = "/?c=hentai&q=%s"

// GrabBP4XType ...
type GrabBP4XType int

// BP4XTypeJAV ...
const (
	BP4XTypeJAV GrabBP4XType = iota
	BP4XTypeAMATEUR
	BP4XTypeIV
	BP4XTypeHENTAI
)

var bp4xGrabList = []string{
	BP4XTypeJAV:     bp4xJavURL,
	BP4XTypeAMATEUR: bp4xAmateurURL,
	BP4XTypeIV:      bp4xIVURL,
	BP4XTypeHENTAI:  bp4xHentaiURL,
}

type grabBP4X struct {
	doc      *goquery.Document
	language GrabLanguage
	grabType GrabBP4XType
	sample   bool
	mainPage string
}

// MainPage ...
func (g *grabBP4X) MainPage(url string) {
	g.mainPage = url
}

// sample ...
func (g *grabBP4X) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabBP4X) Name() string {
	return "bp4x"
}

// Decode ...
func (g *grabBP4X) Decode(*[]*Content) error {
	return nil
}

// Find ...
func (g *grabBP4X) Find(name string) (IGrab, error) {
	name = strings.ToUpper(name)
	url := g.mainPage + bp4xGrabList[g.grabType]
	url = fmt.Sprintf(url, name)
	document, e := net.Query(url)
	if e != nil {
		return g, e
	}
	g.doc = document
	log.Info(g.doc.Text())
	return g, nil
}

// NewGrabBP4X ...
func NewGrabBP4X(grabType GrabBP4XType) IGrab {
	return &grabBP4X{
		mainPage: DefaulBp4xMainPage,
		grabType: grabType,
	}
}
