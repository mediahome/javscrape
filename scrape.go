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
	ExactOff()
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
}

var debug = false

// DefaultInfoName ...
var DefaultInfoName = "inf.json"

// DefaultOutputPath ...
var DefaultOutputPath = "video"

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

// GrabOption ...
func GrabOption(grab IGrab) Options {
	return func(impl *scrapeImpl) {
		grab.SetScrape(impl)
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
	return impl.Range(func(key string, content Content) error {
		e := copyInfo(&content, DefaultOutputPath, strings.ToUpper(key))
		if e != nil {
			return e
		}
		return copyCache(impl.cache, &content, impl.sample, DefaultOutputPath)
	})
}

// ExactOff ...
func (impl *scrapeImpl) ExactOff() {
	for _, g := range impl.grabs {
		g.SetExact(false)
	}
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
		for _, v := range value {
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
		iGrab, e := grab.Find(name)
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "find", name)
			continue
		}
		cs, e := iGrab.Result()
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "decode", name)
		}
		contents = append(contents, cs...)
	}
	if impl.optimize {
		c := MergeOptimize(name, contents)
		e = imageCache(impl.cache, c, impl.sample)
		if e != nil {
			log.Errorw("error", "cache", c.ID, "error", e)
		}
		contents = []*Content{c}
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
	impl.contents[name] = contents
	return nil
}

func (impl *scrapeImpl) init() {
	if impl.cache == nil {
		impl.cache = newCache()
	}
}
