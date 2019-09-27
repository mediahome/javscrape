package net

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Cache ...
type Cache struct {
	tmp string
}

func hash(url string) string {
	sum256 := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", sum256)
}

// NewCache ...
func NewCache(tmp string) *Cache {
	s, e := filepath.Abs(tmp)
	if e != nil {
		panic(e)
	}
	_ = os.MkdirAll(tmp, os.ModePerm)
	return &Cache{tmp: s}
}

// Get ...
func (c *Cache) Get(url string) (e error) {
	_, e = os.Stat(filepath.Join(c.tmp, hash(url)))
	if e == nil || !os.IsNotExist(e) {
		return e
	}
	if cli == nil {
		cli = http.DefaultClient
	}

	res, err := cli.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	name := hash(url)
	file, e := os.OpenFile(filepath.Join(c.tmp, name), os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		return e
	}
	written, e := io.Copy(file, res.Body)
	if e != nil {
		return e
	}
	//ignore written
	_ = written
	return nil
}

// Save ...
func (c *Cache) Save(url string, to string) (written int64, e error) {
	info, e := os.Stat(filepath.Join(c.tmp, hash(url)))
	if e != nil && os.IsNotExist(e) {
		return written, errors.Wrap(e, "cache get error")
	}
	if info.IsDir() {
		return written, errors.New("cache get a dir")
	}
	s, e := filepath.Abs(to)
	if e != nil {
		return written, e
	}
	dir, _ := filepath.Split(s)
	_ = os.MkdirAll(dir, os.ModePerm)
	file, e := os.Open(filepath.Join(c.tmp, hash(url)))
	if e != nil {
		return written, e
	}

	openFile, e := os.OpenFile(filepath.Join(s), os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		return written, e
	}
	return io.Copy(openFile, file)
}
