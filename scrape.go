package scrape

import (
	"strings"
	"sync"

	"github.com/goextension/log"
	"github.com/goextension/log/zap"
)

// RangeFunc ...
type RangeFunc func(key string, content Content) error

// IScrape ...
type IScrape interface {
	Cache() *Cache
	Force(b bool)
	IsGrabSample() (b bool)
	Find(name string) (e error)
	Clear()
	Range(rangeFunc RangeFunc) error
	OutputCallback(f func(key string, content Content) *OutputInfo) []*OutputInfo
	Output() error
}

type scrapeImpl struct {
	contents map[string][]Content
	grabs    []IGrab
	sample   bool
	cache    *Cache
	output   string
	infoName string
	exact    bool
	force    bool
	result   []*OutputInfo
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

func ForceOption(b bool) Options {
	return func(impl *scrapeImpl) {
		impl.force = b
	}
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(opts ...Options) IScrape {
	scrape := &scrapeImpl{
		contents: make(map[string][]Content),
		sample:   true,
		exact:    false,
		output:   DefaultOutputPath,
		infoName: DefaultInfoName,
	}

	for _, opt := range opts {
		opt(scrape)
	}

	scrape.init()

	return scrape
}

func (impl *scrapeImpl) Force(b bool) {
	impl.force = b
}

// Clear ...
func (impl *scrapeImpl) Clear() {
	impl.contents = make(map[string][]Content)
}

// Output ...
func (impl *scrapeImpl) Output() error {
	return impl.Range(func(key string, content Content) (e error) {
		e = copyInfo(&content, DefaultOutputPath, DefaultInfoName)
		if e != nil {
			log.Errorw("copy info", "error", e, "output", key)
		}
		e = copyCache(impl.cache, &content, impl.sample, DefaultOutputPath, false)
		if e != nil {
			log.Errorw("copy cache", "error", e, "output", key)
		}
		return nil
	})
}

func (impl scrapeImpl) OutputCallback(f func(key string, content Content) *OutputInfo) []*OutputInfo {
	var result []*OutputInfo
	impl.Range(func(key string, content Content) error {
		option := f(key, content)
		if option == nil {
			option = DefaultOutputOption()
		}
		if option.Name == "" {
			option.Name = key
		}
		err := copyFileWithInfo(impl.Cache(), content, option)
		if err == nil {
			result = append(result, option)
		}
		return err
	})
	return result
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
			if v.ID == "" {
				continue
			}
			e := rangeFunc(key, v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// Find ...
func (impl *scrapeImpl) Find(name string) (e error) {
	chanContent := make(chan Content, 1)
	wg := &sync.WaitGroup{}
	for _, grab := range impl.grabs {
		wg.Add(1)
		go func(grab IGrab, exact bool, force bool, sample bool, cctx chan<- Content) {
			defer wg.Done()
			grab.SetExact(exact)
			grab.SetSample(sample)
			grab.SetForce(force)
			iGrab, e := grab.Find(name)
			if e != nil {
				log.Errorw("error", "error", e, "name", grab.Name(), "find", name)
				return
			}
			cs, e := iGrab.Result()
			if e != nil {
				log.Errorw("error", "error", e, "name", grab.Name(), "decode", name)
			}
			if debug {
				log.Infow("find", "result", cs)
			}
			for _, c := range cs {
				cctx <- c
			}
		}(grab, impl.exact, impl.force, impl.sample, chanContent)
	}

	go func(cctx chan<- Content) {
		wg.Wait()
		close(cctx)
	}(chanContent)

	for content := range chanContent {
		e = imageCache(impl.cache, content, impl.sample)
		if e != nil {
			log.Errorw("error", "cache", content.ID, "error", e)
		}
		if v, b := impl.contents[content.ID]; b {
			impl.contents[content.ID] = append(v, content)
		} else {
			impl.contents[content.ID] = []Content{content}
		}
		if debug {
			log.Infow("find", "content", content)
		}
	}

	return nil
}

func (impl *scrapeImpl) init() {
	if impl.cache == nil {
		impl.cache = NewCache()
	}
}
