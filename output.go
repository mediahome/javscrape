package scrape

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goextension/log"
)

// DefaultOutputPath ...
var DefaultOutputPath = "image"

type OutputInfo struct {
	Skip       bool
	Force      bool
	OutputPath string
	CopyInfo   bool
	InfoPath   string
	InfoName   string
	CopyPoster bool
	PosterPath string
	PosterName string
	CopyThumb  bool
	ThumbPath  string
	ThumbName  string
	CopySample bool
	SamplePath string
	SampleName string
	ImagePath  string
}

func DefaultOutputOption() *OutputInfo {
	return &OutputInfo{
		Skip:       false,
		OutputPath: "",
		ImagePath:  "image",
		CopyInfo:   false,
		InfoPath:   "",
		InfoName:   ".nfo",
		CopyPoster: true,
		PosterPath: "",
		PosterName: "poster",
		CopyThumb:  true,
		ThumbPath:  "",
		ThumbName:  "thumb",
		CopySample: false,
		SamplePath: "",
		SampleName: "sample",
	}
}

func copyCache(cache *Cache, msg *Content, sample bool, output string, force bool) (e error) {
	pid := filepath.Join(output, strings.ToUpper(msg.ID), "."+msg.From)
	e = copyFile(cache, msg.Poster, filepath.Join(pid, "poster"), force)
	if e != nil {
		return e
	}
	e = copyFile(cache, msg.Thumb, filepath.Join(pid, "thumb"), force)
	if e != nil {
		return e
	}
	for _, act := range msg.Actors {
		e = copyFile(cache, act.Image, filepath.Join(pid, ".actor", act.Name), force)
		if e != nil {
			return e
		}
	}
	if sample {
		for _, s := range msg.Sample {
			e = copyFile(cache, s.Image, filepath.Join(pid, ".sample", "sample"+"@"+strconv.Itoa(s.Index)), force)
			if e != nil {
				return e
			}
			e = copyFile(cache, s.Thumb, filepath.Join(pid, ".thumb", "thumb"+"@"+strconv.Itoa(s.Index)), force)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func copyInfo(msg *Content, path string, name string) error {
	inf := filepath.Join(path, name)
	_ = os.MkdirAll(filepath.Dir(inf), os.ModePerm)
	info, e := os.Stat(inf)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}
	bytes, e := json.MarshalIndent(msg, "", " ")
	if e != nil {
		return e
	}
	return ioutil.WriteFile(inf, bytes, 0755)
}

func copyFileWithInfo(cache *Cache, content Content, option *OutputInfo) error {
	var e error
	if option.Skip {
		return nil
	}
	if option.CopyInfo {
		e = copyInfo(&content, filepath.Join(option.OutputPath, option.InfoPath), option.InfoName)
		if e != nil {
			log.Errorw("OutputCallback", "error", e, "output", content.ID)
		}
	}

	if option.CopyPoster {
		ext := Ext(content.Poster)
		option.PosterName = option.PosterName + ext
		path := filepath.Join(option.OutputPath, option.ImagePath, option.PosterPath, option.PosterName)
		if debug {
			log.Infow("CopyFile", "source", content.Poster, "path", path)
		}
		e = copyFile(cache, content.Poster, path, option.Force)
		if e != nil {
			log.Errorw("OutputCallback", "error", e, "output", content.ID)
		}
	}

	if option.CopyThumb {
		ext := Ext(content.Thumb)
		option.ThumbName = option.ThumbName + ext
		path := filepath.Join(option.OutputPath, option.ImagePath, option.ThumbPath, option.ThumbName)
		e = copyFile(cache, content.Thumb, path, option.Force)
		if e != nil {
			log.Errorw("OutputCallback", "error", e, "output", content.ID)
		}
	}

	if option.CopySample {
		for i, sample := range content.Sample {
			ext := Ext(content.Thumb)
			option.ThumbName = option.ThumbName + "@" + strconv.Itoa(i) + ext
			path := filepath.Join(option.OutputPath, option.ImagePath, option.SamplePath, option.SampleName)
			e = copyFile(cache, sample.Image, path, option.Force)
			if e != nil {
				log.Errorw("OutputCallback", "error", e, "output", content.ID)
			}
		}
	}
	return e
}

func copyFile(cache *Cache, source, path string, force bool) error {
	if source == "" {
		return nil
	}

	if debug {
		log.Infow("CopyFile", "source", source, "path", path)
	}
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	info, e := os.Stat(path)
	if e != nil && !os.IsNotExist(e) {
		return e
	}

	var bys []byte
	if e == nil && info.Size() != 0 {
		return nil
	}
	if force {
		bys, e = cache.ForceGet(source)
	} else {
		bys, e = cache.GetBytes(source)
	}
	if e != nil {
		return e
	}
	return ioutil.WriteFile(path, bys, 0755)
}

func imageCache(cache *Cache, m Content, sample bool) (e error) {
	path := make(chan string, 2)
	go func(path chan<- string) {
		defer close(path)
		path <- m.Poster
		path <- m.Thumb
		for _, act := range m.Actors {
			path <- act.Image
		}
		if sample {
			for _, s := range m.Sample {
				path <- s.Image
				path <- s.Thumb
			}
		}
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
