package network

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/goextension/log"
	"golang.org/x/net/proxy"
)

//var log = trait.NewZapSugar()
var (
	cli *http.Client
)

func init() {
	cli = http.DefaultClient
}

// ProxyArgs ...
type ProxyArgs func(cli *http.Client)

// TimeOut ...
func TimeOut(sec int) ProxyArgs {
	return func(cli *http.Client) {
		cli.Timeout = time.Duration(sec) * time.Second
	}
}

// RegisterProxy ...
func RegisterProxy(addr string, args ...ProxyArgs) (e error) {
	log.Debug("NETWORK", "register proxy", "addr", addr)
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

	for _, fn := range args {
		fn(cli)
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

func Client() *http.Client {
	return cli
}
