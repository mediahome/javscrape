package net

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
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
		Timeout:       60 * time.Second,
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
