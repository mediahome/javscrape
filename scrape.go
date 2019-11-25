package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goextension/log"
)

// IScrape ...
type IScrape interface {
	IsGrabSample() (b bool)
	Find(name string) (msg *[]*Content, e error)
}

type scrapeImpl struct {
	contents []*Content
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
		//grabs: grabs,
		//sample:   false,
		//cache:    nil,
		output:   DefaultOutputPath,
		infoName: DefaultInfoName,
	}

	for _, opt := range opts {
		opt(scrape)
	}
	return scrape
}

// Find ...
func (impl *scrapeImpl) Find(name string) (msg *[]*Content, e error) {
	for _, grab := range impl.grabs {
		var c Content
		iGrab, e := grab.Find(name)
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "find", name)
			continue
		}
		e = iGrab.Decode(&c)
		if e != nil {
			log.Errorw("error", "error", e, "name", grab.Name(), "decode", name)
		}
		impl.contents = append(impl.contents, &c)
	}

	if len(impl.contents) == 0 {
		return nil, fmt.Errorf("[%s] not found", name)
	}

	if impl.cache != nil {
		for _, m := range impl.contents {
			e := imageCache(impl.cache, m)
			if e != nil {
				return nil, e
			}
		}
	}

	var err error
	if impl.output != "" {
		for _, m := range impl.contents {
			e = copyInfo(m, impl.output, impl.infoName)
			if e != nil {
				log.Errorw("error", "error1", e, "msg", m)
				err = e
			}
			e = copyCache(impl.cache, m, impl.output)
			if e != nil {
				log.Errorw("error", "error2", e, "msg", m)
				err = e
			}
		}
	}
	return msg, err
}

func copyCache(cache *Cache, msg *Content, output string) (e error) {
	pid := filepath.Join(output, strings.ToUpper(msg.ID))
	e = copyFile(cache, msg.Image, filepath.Join(pid, "image"))
	if e != nil {
		return e
	}
	e = copyFile(cache, msg.Thumb, filepath.Join(pid, "thumb"))
	if e != nil {
		return e
	}
	for _, act := range msg.Actors {
		e = copyFile(cache, act.Image, filepath.Join(pid, ".actor", act.Name))
		if e != nil {
			return e
		}
	}
	for _, s := range msg.Sample {
		e = copyFile(cache, s.Image, filepath.Join(pid, ".sample", "sample"+"@"+strconv.Itoa(s.Index)))
		if e != nil {
			return e
		}
		e = copyFile(cache, s.Thumb, filepath.Join(pid, ".thumb", "thumb"+"@"+strconv.Itoa(s.Index)))
		if e != nil {
			return e
		}
	}
	return nil
}

func copyInfo(msg *Content, path string, name string) error {
	pid := filepath.Join(path, strings.ToUpper(msg.ID))
	inf := filepath.Join(pid, name)
	_ = os.MkdirAll(filepath.Dir(inf), os.ModePerm)
	info, e := os.Stat(inf)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}
	file, e := os.OpenFile(inf, os.O_SYNC|os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")
	return enc.Encode(msg)
}

// TrimEnd ...
func TrimEnd(source string) string {
	return strings.Split(source, "?")[0]
}

// Ext ...
func Ext(source string) string {
	ext := filepath.Ext(TrimEnd(source))
	if debug {
		log.Infow("ext", "source", source, "ext", ext)
	}
	return ext
}

func copyFile(cache *Cache, source, path string) error {
	if source == "" {
		return nil
	}
	reader, e := cache.Reader(source)
	if e != nil {
		return e
	}
	path = TrimEnd(path)
	if debug {
		log.Infow("copy", "dir", filepath.Dir(path), "path", path)
	}
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	info, e := os.Stat(path + Ext(source))
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}

	file, e := os.OpenFile(path+Ext(source), os.O_SYNC|os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer file.Close()
	written, e := io.Copy(file, reader)
	if e != nil {
		return e
	}
	_ = written
	return nil
}

func imageCache(cache *Cache, m *Content) (e error) {
	path := make(chan string)
	go func(path chan<- string) {
		defer close(path)
		//for _, m := range msg {
		path <- m.Image
		path <- m.Thumb
		for _, act := range m.Actors {
			path <- act.Image
		}
		for _, s := range m.Sample {
			path <- s.Image
			path <- s.Thumb
		}
		//}
	}(path)

	for p := range path {
		if p != "" {
			_, err := cache.Get(p)
			if err != nil && !os.IsExist(err) {
				log.Error(err)
			}
		}
	}
	return nil
}
