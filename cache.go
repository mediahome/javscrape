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

	"github.com/gocacher/cacher"
	"github.com/gocacher/file-cache"
	"github.com/goextension/log"
)

// DefaultCachePath ...
var DefaultCachePath = "tmp"

// Cache ...
type Cache struct {
	cache cacher.Cacher
}

var _cache *Cache

// HasCache ...
var HasCache bool

func init() {
	HasCache = true
	_cache = newCache()
}

func newCache() *Cache {
	cache.DefaultPath = DefaultCachePath
	return &Cache{
		cache: &cache.FileCache{},
	}

}

// Hash ...
func Hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", sum256)
}

// Reader ...
func (c *Cache) Reader(url string) (io.Reader, error) {
	return c.Get(url)
}

// Get ...
func (c *Cache) Get(url string) (reader io.Reader, e error) {
	name := Hash(url)
	log.Infow("cache get", "url", url, "hash", name)
	b, e := c.cache.Has(name)
	if e == nil && b {
		getted, e := c.cache.Get(name)
		if e != nil {
			return nil, e
		}
		return bytes.NewReader(getted), nil
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
	bys, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	e = c.cache.Set(name, bys)
	if e != nil {
		return nil, e
	}
	return bytes.NewReader(bys), nil
}

// Save ...
func (c *Cache) Save(path, url, to string) (written int64, e error) {
	s, e := filepath.Abs(to)
	if e != nil {
		return written, e
	}
	dir, _ := filepath.Split(s)
	_ = os.MkdirAll(dir, os.ModePerm)
	fromData, e := c.cache.Get(Hash(url))
	if e != nil {
		return 0, e
	}

	toFile, e := os.OpenFile(s, os.O_TRUNC|os.O_CREATE|os.O_RDWR|os.O_SYNC, os.ModePerm)
	if e != nil {
		return written, e
	}

	n, e := toFile.Write(fromData)
	if e != nil {
		return 0, e
	}
	return int64(n), nil
}
