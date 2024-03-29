package core

import (
	"net/url"
	"path"
)

func URL(prefix string, uris ...string) string {
	u, err := url.Parse(prefix)
	if err != nil {
		return prefix
	}
	uris = append([]string{u.Path}, uris...)
	u.Path = path.Join(uris...)
	return u.String()
}

func URLAddValues(urlpath string, v url.Values) string {
	if v == nil {
		return urlpath
	}
	return urlpath + "?" + v.Encode()
}
