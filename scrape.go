package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/javscrape/go-scrape/net"
)

var debug = false

// IScrape ...
type IScrape interface {
	GrabSample(b bool)
	IsGrabSample() (b bool)
	ImageCache(path string)
	Output(path string)
	Find(name string) (msg *[]*Content, e error)
}

type scrapeImpl struct {
	grabs  []IGrab
	sample bool
	//cache  string
	cache  *net.Cache
	output string
}

// Output ...
func (impl *scrapeImpl) Output(path string) {
	impl.output = path
}

// ImageCache ...
func (impl *scrapeImpl) ImageCache(path string) {
	impl.cache = net.NewCache(path)
}

// IsGrabSample ...
func (impl *scrapeImpl) IsGrabSample() bool {
	return impl.sample
}

// GrabSample ...
func (impl *scrapeImpl) GrabSample(b bool) {
	impl.sample = b
	if !impl.sample {
		return
	}
	for _, grab := range impl.grabs {
		grab.Sample(b)
	}
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(grabs ...IGrab) IScrape {
	return &scrapeImpl{grabs: grabs}
}

// Find ...
func (impl *scrapeImpl) Find(name string) (msg *[]*Content, e error) {
	msg = new([]*Content)
	for _, grab := range impl.grabs {
		iGrab, e := grab.Find(name)
		if e != nil {
			log.With("name", grab.Name(), "find", name).Error(e)
			continue
		}
		e = iGrab.Decode(msg)
		if e != nil {
			log.With("name", grab.Name(), "decode", name).Error(e)
		}
	}

	if len(*msg) == 0 {
		return nil, fmt.Errorf("[%s] not found", name)
	}

	if impl.cache != nil {
		e := imageCache(impl.cache, *msg)
		if e != nil {
			return nil, e
		}
	}

	var err error
	if impl.output != "" {
		for _, m := range *msg {
			e = copyInfo(m, impl.output)
			if e != nil {
				log.With("msg1", m).Error(e)
				err = e
			}
			e = copyCache(impl.cache, m, impl.output)
			if e != nil {
				log.With("msg2", m).Error(e)
				err = e
			}
		}
	}
	return msg, err
}

func copyCache(cache *net.Cache, msg *Content, output string) (e error) {
	pid := filepath.Join(output, msg.ID)
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
		e = copyFile(cache, s.Image, filepath.Join(pid, ".sample", "image"+"_"+strconv.Itoa(s.Index)))
		if e != nil {
			return e
		}
		e = copyFile(cache, s.Thumb, filepath.Join(pid, ".sample", "thumb"+"_", strconv.Itoa(s.Index)))
		if e != nil {
			return e
		}
	}
	return nil
}

func copyInfo(msg *Content, path string) error {
	pid := filepath.Join(path, msg.ID)
	inf := filepath.Join(pid, "inf.json")
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
	return enc.Encode(msg)
}

func copyFile(cache *net.Cache, source, path string) error {
	if source == "" {
		return nil
	}
	reader, e := cache.Reader(source)
	if e != nil {
		return e
	}
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	ext := filepath.Ext(source)
	info, e := os.Stat(path + ext)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}

	file, e := os.OpenFile(path+ext, os.O_SYNC|os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
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

func imageCache(cache *net.Cache, msg []*Content) (e error) {
	path := make(chan string)
	go func(path chan<- string) {
		defer close(path)
		for _, m := range msg {
			path <- m.Image
			path <- m.Thumb
			for _, act := range m.Actors {
				path <- act.Image
			}
			for _, s := range m.Sample {
				path <- s.Image
				path <- s.Thumb
			}
		}
	}(path)

	for p := range path {
		if p != "" {
			err := cache.Get(p)
			if err != nil && !os.IsExist(err) {
				log.Error(err)
			}
		}
	}
	return nil
}
