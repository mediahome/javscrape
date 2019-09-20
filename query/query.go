package query

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

var queryProxy proxy.Dialer

// ProxyURL ...
var ProxyURL string

// RegisterProxy ...
func RegisterProxy(path string) {
	proxyURL, err := url.Parse(path)
	if err != nil {
		return
	}
	p, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return
	}

	host := proxy.NewPerHost(p, proxy.Direct)
	host.AddFromString("localhost, 127.0.0.1")
}

func getTransport() *http.Transport {
	if ProxyURL != "" {
		proxy, _ := url.Parse(ProxyURL)
		return &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &http.Transport{
		Proxy:           nil,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// New ...
func New(url string) (*goquery.Document, error) {
	cli := &http.Client{
		Transport:     getTransport(),
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       15 * time.Second,
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
