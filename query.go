package scrape

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"
)

// Query ...
func Query(url string) (*goquery.Document, error) {
	if HasCache {
		closer, e := _cache.Reader(url)
		if e != nil {
			return nil, e
		}
		return goquery.NewDocumentFromReader(closer)
	}
	if cli == nil {
		cli = http.DefaultClient
	}

	res, err := cli.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	return doc, err
}
