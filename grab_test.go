package scrape

import (
	"testing"

	"github.com/javscrape/go-scrape/net"
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
	msg := new([]*Content)
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
	doc, err := grab.Find("abp/10")
	msg := new([]*Content)
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
