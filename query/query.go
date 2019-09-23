package query

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

var queryProxy proxy.Dialer

// RegisterProxy ...
func RegisterProxy(addr string) (e error) {
	u, e := url.Parse(addr)
	if e != nil {
		return e
	}
	switch u.Scheme {
	case "socks5":
		queryProxy, e = proxySOCKS5(u.Host)
	}
	return nil
}

func proxySOCKS5(addr string) (proxy.Dialer, error) {
	return proxy.SOCKS5("tcp", addr,
		nil, //&proxy.Auth{User: "", Password: ""},
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)
}

func getTransport() *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			if queryProxy != nil {
				return queryProxy.Dial(network, addr)
			}
			return net.Dial(network, addr)
		},
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
