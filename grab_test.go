package scrape

import (
	"github.com/javscrape/go-scrape/net"
	"testing"
)

// TestGrabBP4X_Find ...
func TestGrabBP4X_Find(t *testing.T) {
	e := net.RegisterProxy("socks5://localhost:11080")
	if e != nil {
		return
	}
	grab := NewGrabBP4X(BP4XTypeJAV)
	doc, err := grab.Find("abp-874")
	if err != nil {
		t.Fatal(err)
	}
	msg := new([]*Message)
	err = doc.Decode(msg)
	if err != nil {
		t.Fatal(err)
	}
	cache := net.NewCache("./tmp")

	err = imageCache(cache, *msg)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGrabJAVBUS_Find ...
func TestGrabJAVBUS_Find(t *testing.T) {
	DebugOn()
	e := net.RegisterProxy("socks5://localhost:10808")
	if e != nil {
		return
	}
	grab := NewGrabJAVBUS()
	grab.Sample(true)
	doc, err := grab.Find("gah-11")
	msg := new([]*Message)
	if err != nil {
		t.Fatal(err)
	}
	err = doc.Decode(msg)
	t.Logf("%+v", *msg)
	if err != nil {
		t.Fatal(err)
	}
	cache := net.NewCache("./tmp")

	err = imageCache(cache, *msg)
	if err != nil {
		t.Fatal(err)
	}
}
