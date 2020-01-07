package scrape

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"github.com/goextension/log"
)

// DefaultCachePath ...
var DefaultCachePath = "tmp"

// Cache ...
type Cache struct {
	lock  sync.Mutex
	cache cacher.Cacher
}

func newCache() *Cache {
	cache.DefaultCachePath = DefaultCachePath
	return &Cache{
		cache: cache.New(),
	}
}

// Hash ...
func Hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", sum256)
}

// GetReader ...
func (c *Cache) GetReader(url string) (io.Reader, error) {
	bys, e := c.Get(url)
	if e != nil {
		return nil, e
	}
	return bytes.NewReader(bys), nil
}

// GetBytes ...
func (c *Cache) GetBytes(url string) ([]byte, error) {
	return c.Get(url)
}

// Get ...
func (c *Cache) Get(url string) (bys []byte, e error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	name := Hash(url)
	log.Infow("cache get", "url", url, "hash", name)
	b, e := c.cache.Has(name)
	if e == nil && b {

		getted, e := c.cache.Get(name)
		if e != nil {
			return nil, e
		}
		return getted, nil
	}

	if cli == nil {
		cli = http.DefaultClient
	}

	res, e := cli.Get(url)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	bys, e = ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	e = c.cache.Set(name, bys)
	if e != nil {
		return nil, e
	}
	return bys, nil
}

// Save ...
func (c *Cache) Save(url, to string) (e error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	s, e := filepath.Abs(to)
	if e != nil {
		return e
	}
	dir, _ := filepath.Split(s)
	_ = os.MkdirAll(dir, os.ModePerm)
	fromData, e := c.Get(url)
	if e != nil {
		return e
	}

	e = ioutil.WriteFile(s, fromData, 0755)
	if e != nil {
		return e
	}
	return nil
}

// Query ...
func (c *Cache) Query(url string) (*goquery.Document, error) {
	closer, e := c.GetReader(url)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}
