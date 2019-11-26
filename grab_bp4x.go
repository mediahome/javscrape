package scrape

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"
)

// DefaulBp4xMainPage ...
const DefaulBp4xMainPage = "https://www.bp4x.com"
const bp4xJavURL = "/?q=%s"
const bp4xAmateurURL = "/?c=amateur&q=%s"
const bp4xIVURL = "/?c=iv&q=%s"
const bp4xHentaiURL = "/?c=hentai&q=%s"

// GrabBp4xType ...
type GrabBp4xType int

// BP4XTypeJAV ...
const (
	BP4XTypeJAV GrabBp4xType = iota
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

type grabBp4x struct {
	scrape   IScrape
	doc      *goquery.Document
	language GrabLanguage
	grabType GrabBp4xType
	sample   bool
	mainPage string
}

// SetExact ...
func (g *grabBp4x) SetExact(bool) {
	panic("implement me")
}

// Result ...
func (g *grabBp4x) Result() ([]*Content, error) {
	panic("implement me")
}

// SetSample ...
func (g *grabBp4x) SetSample(bool) {
	panic("implement me")
}

// SetScrape ...
func (g *grabBp4x) SetScrape(scrape IScrape) {
	g.scrape = scrape
}

// Clone ...
func (g *grabBp4x) Clone() IGrab {
	panic("implement me")
}

// HasNext ...
func (g *grabBp4x) HasNext() bool {
	panic("implement me")
}

// Next ...
func (g *grabBp4x) Next() (IGrab, error) {
	panic("implement me")
}

// MainPage ...
func (g *grabBp4x) MainPage(url string) {
	g.mainPage = url
}

// sample ...
func (g *grabBp4x) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabBp4x) Name() string {
	return "bp4x"
}

// Decode ...
func (g *grabBp4x) Decode(*Content) error {
	return nil
}

// Find ...
func (g *grabBp4x) Find(name string) (IGrab, error) {
	name = strings.ToUpper(name)
	url := g.mainPage + bp4xGrabList[g.grabType]
	url = fmt.Sprintf(url, name)
	document, e := g.scrape.Cache().Query(url)
	if e != nil {
		return g, e
	}
	g.doc = document
	log.Info(g.doc.Text())
	return g, nil
}

// GrabBp4xOptions ...
type GrabBp4xOptions func(javbus *grabBp4x)

// GrabBp4xTypeOption ...
func GrabBp4xTypeOption(grabType GrabBp4xType) GrabBp4xOptions {
	return func(grab *grabBp4x) {
		grab.grabType = grabType
	}
}

// NewGrabBp4x ...
func NewGrabBp4x(ops ...GrabBp4xOptions) IGrab {
	grab := &grabBp4x{
		mainPage: DefaulBp4xMainPage,
		grabType: BP4XTypeJAV,
	}
	for _, op := range ops {
		op(grab)
	}
	return grab
}
