package scrape

import (
	"fmt"
	"strings"

	"github.com/goextension/log"
	"github.com/goextension/log/zap"
)

// RangeFunc ...
type RangeFunc func(key string, content Content) error

// IScrape ...
type IScrape interface {
	Cache() *Cache
	IsGrabSample() (b bool)
	Find(name string) (e error)
	Clear()
	Range(rangeFunc RangeFunc) error
	Output() error
}

type scrapeImpl struct {
	contents map[string][]*Content
	grabs    []IGrab
	sample   bool
	cache    *Cache
	output   string
	infoName string
	optimize bool
	exact    bool
}

var debug = false

// DefaultInfoName ...
var DefaultInfoName = ".info"

// Options ...
type Options func(impl *scrapeImpl)

// IsGrabSample ...
func (impl *scrapeImpl) IsGrabSample() bool {
	return impl.sample
}

// CacheOption ...
func CacheOption(cache *Cache) Options {
	return func(impl *scrapeImpl) {
		impl.cache = cache
	}
}

// OptimizeOption ...
func OptimizeOption(b bool) Options {
	return func(impl *scrapeImpl) {
		impl.optimize = b
	}
}

// MergeOptimize ...
func MergeOptimize(id string, contents []*Content) *Content {
	var content *Content
	for _, c := range contents {
		if strings.ToUpper(c.ID) == strings.ToUpper(id) {
			if content == nil {
				content = c
				continue
			}
		}
		if content == nil {
			continue
		}
		if strings.ToUpper(content.ID) == strings.ToUpper(c.ID) {
			if len(content.Sample) < len(c.Sample) {
				log.Infow("optimize", "field", "sample")
				content.Sample = c.Sample
			}
			if len(content.Genres) < len(c.Genres) {
				log.Infow("optimize", "field", "genre")
				content.Genres = c.Genres
			}
			if len(content.Actors) < len(c.Actors) {
				log.Infow("optimize", "field", "actor")
				content.Actors = c.Actors
			}
		}
	}
	return content
}

// SampleOption ...
func SampleOption(b bool) Options {
	return func(impl *scrapeImpl) {
		impl.sample = b
	}
}

// ExactOption ...
func ExactOption(b bool) Options {
	return func(impl *scrapeImpl) {
		impl.exact = b
	}
}

// GrabOption ...
func GrabOption(grab IGrab) Options {
	return func(impl *scrapeImpl) {
		impl.grabs = append(impl.grabs, grab)
	}
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(opts ...Options) IScrape {
	scrape := &scrapeImpl{
		contents: make(map[string][]*Content),
		sample:   true,
		exact:    false,
		optimize: true,
		output:   DefaultOutputPath,
		infoName: DefaultInfoName,
	}

	for _, opt := range opts {
		opt(scrape)
	}

	scrape.init()

	return scrape
}

// Clear ...
func (impl *scrapeImpl) Clear() {
	impl.contents = make(map[string][]*Content)
}

// Output ...
func (impl *scrapeImpl) Output() error {
	return impl.Range(func(key string, content Content) (e error) {
		e = copyInfo(&content, DefaultOutputPath, DefaultInfoName)
		if e != nil {
			log.Errorw("copy info", "error", e, "output", key)
		}
		e = copyCache(impl.cache, &content, impl.sample, DefaultOutputPath)
		if e != nil {
			log.Errorw("copy cache", "error", e, "output", key)
		}
		return nil
	})
}

func init() {
	zap.InitZapSugar()
}

// Cache ...
func (impl *scrapeImpl) Cache() *Cache {
	return impl.cache
}

// Range ...
func (impl *scrapeImpl) Range(rangeFunc RangeFunc) error {
	for key, value := range impl.contents {
		log.Infow("range", "key", key)
		for _, v := range value {
			if v == nil {
				continue
			}
			e := rangeFunc(key, *v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// Find ...
func (impl *scrapeImpl) Find(name string) (e error) {
	var contents []*Content
	for _, grab := range impl.grabs {
		grab.SetExact(impl.exact)
		grab.SetSample(impl.sample)
		iGrab, e := grab.Find(name)
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "find", name)
			continue
		}
		cs, e := iGrab.Result()
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "decode", name)
		}
		if debug {
			log.Infow("find", "result", cs)
		}
		contents = append(contents, cs...)
	}
	if impl.exact && impl.optimize {
		c := MergeOptimize(name, contents)
		if c != nil {
			e = imageCache(impl.cache, c, impl.sample)
			if e != nil {
				log.Errorw("error", "cache", c.ID, "error", e)
			}
			contents = []*Content{c}
		}
	} else {
		for _, c := range contents {
			e = imageCache(impl.cache, c, impl.sample)
			if e != nil {
				log.Errorw("error", "cache", c.ID, "error", e)
			}
		}
	}
	if len(contents) == 0 {
		return fmt.Errorf("[%s] not found", name)
	}
	if debug {
		log.Infow("find", "contents", contents)
	}
	impl.contents[name] = contents
	return nil
}

func (impl *scrapeImpl) init() {
	if impl.cache == nil {
		impl.cache = NewCache()
	}
}
