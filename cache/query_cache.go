package cache

import (
	"github.com/PuerkitoBio/goquery"
)

type QueryCache NetCache

// Query ...
func (c *QueryCache) Query(url string, force bool) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *QueryCache) GetQuery(url string, force bool) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *QueryCache) ForceQuery(url string) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, true)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *QueryCache) net() *NetCache {
	return (*NetCache)(c)
}

func NewQueryCache() *QueryCache {
	return (*QueryCache)(newCache())
}

var _ = (*QueryCache)(&NetCache{})
