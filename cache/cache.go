package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// DefaultCachePath ...
var DefaultCachePath = "tmp"
var _cache NetCacher
var _cacheOnce *sync.Once

// Querier ...
// @Description:
type Querier interface {
	Query(url string, force bool) (*goquery.Document, error)
	GetQuery(url string, force bool) (*goquery.Document, error)
	ForceQuery(url string) (*goquery.Document, error)
}

type NetCacher interface {
	GetReader(url string, force bool) (io.Reader, error)
	GetBytes(url string, force bool) ([]byte, error)
	HasURL(url string) bool
	Save(url, to string) (e error)
}

// netCache ...

func init() {
	_cacheOnce = &sync.Once{}
}

// Hash ...
func Hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return hex.EncodeToString(sum256[:])
}
