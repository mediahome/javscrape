package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/PuerkitoBio/goquery"
	cache "github.com/gocacher/badger-cache/v3"
	"github.com/gocacher/cacher"

	"github.com/javscrape/go-scrape/network"
)

// DefaultCachePath ...
var DefaultCachePath = "tmp"
var _cache *Cache
var _cacheOnce *sync.Once

// Cache ...
type Cache struct {
	lock sync.Mutex
	cacher.Cacher
}

func init() {
	_cacheOnce = &sync.Once{}
}

func newCache() *Cache {
	cache.DefaultCachePath = DefaultCachePath
	return &Cache{
		Cacher: cache.New(),
	}
}

// NewCache ...
func NewCache() *Cache {
	_cacheOnce.Do(func() {
		_cache = newCache()
	})
	return _cache
}

// Hash ...
func Hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return hex.EncodeToString(sum256[:])
}

// GetReader ...
func (c *Cache) GetReader(url string, force bool) (io.Reader, error) {
	bys, e := c.get(url, force)
	if e != nil {
		return nil, e
	}
	return bytes.NewReader(bys), nil
}

// GetBytes ...
func (c *Cache) GetBytes(url string, force bool) ([]byte, error) {
	return c.get(url, force)
}

func (c *Cache) HasURL(url string) bool {
	return c.has(Hash(url))
}

func (c *Cache) has(name string) bool {
	exist, err := c.Has(name)
	return err == nil && exist
}

func (c *Cache) get(url string, useCache bool) (bys []byte, e error) {
	name := Hash(url)
	if useCache {
		b := c.has(name)
		//log.Infow("cache get", "url", url, "hash", name, "exist", b)
		if b {
			bys, e = c.Get(name)
			if e != nil {
				return nil, e
			}
			return bys, nil
		}
	}
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
	bys, e = ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	e = c.Set(name, bys)
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
func (c *Cache) Query(url string, force bool) (*goquery.Document, error) {
	closer, e := c.GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *Cache) GetQuery(url string, force bool) (*goquery.Document, error) {
	closer, e := c.GetReader(url, force)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}

func (c *Cache) ForceQuery(url string) (*goquery.Document, error) {
	closer, e := c.GetReader(url, true)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(closer)
}
