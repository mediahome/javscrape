package scrape

import (
	"github.com/javscrape/go-scrape/query"
	"testing"
)

// TestGrabBP4X_Find ...
func TestGrabBP4X_Find(t *testing.T) {
	grab := NewGrabBP4X(BP4XTypeJAV)
	doc, err := grab.Find("abp-874")
	msg := *new([]*Message)
	err = doc.Decode(msg)
	t.Log(err)
}

// TestGrabJAVBUS_Find ...
func TestGrabJAVBUS_Find(t *testing.T) {
	DebugOn()
	e := query.RegisterProxy("socks5://localhost:1080")
	if e != nil {
		return
	}
	grab := NewGrabJAVBUS(LanguageEnglish)
	doc, err := grab.Find("abp-773")
	msg := *new([]*Message)
	if err != nil {
		t.Fatal(err)
	}
	err = doc.Decode(msg)
	t.Logf("%+v", msg)
	t.Error(err)
}
