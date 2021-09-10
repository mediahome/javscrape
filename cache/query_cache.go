package cache

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocacher/cacher"
)

type queryCache netCache

func (c *queryCache) Cache() cacher.Cacher {
	return c.net().Cache()
}

// Query ...
func (c *queryCache) Query(url string, force bool) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *queryCache) GetQuery(url string, force bool) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *queryCache) ForceQuery(url string) (*goquery.Document, error) {
	closer, e := c.net().GetReader(url, true)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *queryCache) net() *netCache {
	return (*netCache)(c)
}

func NewQueryCache(client *http.Client) Querier {
	return (*queryCache)(newCache(client))
}

var _ = Querier((*queryCache)(&netCache{}))
var _ = (*queryCache)(&netCache{})
