package cache

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/javscrape/go-scrape/network"
)

type nocache struct {
	client *http.Client
}

func (n nocache) getReader(url string) (io.Reader, error) {
	cli := network.Client()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.11 Safari/537.36")

	res, e := cli.Do(req)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	bys, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	return bytes.NewReader(bys), nil
}

func (n nocache) Query(url string, force bool) (*goquery.Document, error) {
	reader, err := n.getReader(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func (n nocache) GetQuery(url string, force bool) (*goquery.Document, error) {
	reader, err := n.getReader(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func (n nocache) ForceQuery(url string) (*goquery.Document, error) {
	reader, err := n.getReader(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func NoCacheQuery(client *http.Client) Querier {
	return &nocache{
		client: client,
	}
}
