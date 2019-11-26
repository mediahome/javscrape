package scrape

import (
	"fmt"

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
}

type scrapeImpl struct {
	contents map[string][]*Content
	grabs    []IGrab
	sample   bool
	cache    *Cache
	output   string
	infoName string
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

	if len(contents) == 0 {
		return fmt.Errorf("[%s] not found", name)
	}
	impl.contents[name] = contents
	return nil
	//if impl.cache != nil {
	//	for _, m := range impl.contents {
	//		e := imageCache(impl.cache, m)
	//		if e != nil {
	//			return nil, e
	//		}
	//	}
	//}
	//
	//var err error
	//if impl.output != "" {
	//	for _, m := range impl.contents {
	//		e = copyInfo(m, impl.output, impl.infoName)
	//		if e != nil {
	//			log.Errorw("error", "error1", e, "msg", m)
	//			err = e
	//		}
	//		e = copyCache(impl.cache, m, impl.output)
	//		if e != nil {
	//			log.Errorw("error", "error2", e, "msg", m)
	//			err = e
	//		}
	//	}
	//}
	//return msg, err
}

func (impl *scrapeImpl) init() {
	if impl.cache == nil {
		impl.cache = newCache()
	}
}
