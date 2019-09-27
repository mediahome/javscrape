package net

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

var cli *http.Client

// RegisterProxy ...
func RegisterProxy(addr string) (e error) {
	u, e := url.Parse(addr)
	if e != nil {
		return e
	}
	var transport *http.Transport
	switch u.Scheme {
	case "http", "https":
		transport = getHTTPTransport(u)
	case "socks5":
		transport = getSOCKS5Transport(u.Host)
	default:
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	cli = &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       15 * time.Second,
	}
	return nil
}

func getHTTPTransport(u *url.URL) *http.Transport {
	return &http.Transport{
		Proxy:           http.ProxyURL(u),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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
func getSOCKS5Transport(addr string) *http.Transport {
	queryProxy, err := proxySOCKS5(addr)
	if err != nil {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return queryProxy.Dial(network, addr)
		},
	}
}

// NewQuery ...
func NewQuery(url string) (*goquery.Document, error) {
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
