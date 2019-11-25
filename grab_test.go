package scrape

import (
	"os"
	"testing"
)

// TestGrabJAVBUS_Find ...
func TestGrabJAVBUS_Find(t *testing.T) {
	DebugOn()
	e := RegisterProxy("socks5://localhost:10808")
	if e != nil {
		return
	}
	grab := NewGrabJavbus()
	grab.SetSample(true)
	doc, err := grab.Find("abp-874")

	if err != nil {
		t.Fatal(err)
	}
	var msg Content
	err = doc.Decode(&msg)
	t.Logf("%+v", msg)
	if err != nil {
		t.Fatal(err)
	}
	cache := newCache()

	err = imageCache(cache, &msg)
	if err != nil {
		t.Fatal(err)
	}
}

// TestStat ...
func TestStat(t *testing.T) {
	info, e := os.Stat("grab_test.go")
	t.Log(os.IsNotExist(e))
	t.Log(os.IsExist(e))
	t.Log(info, e)

	info1, e1 := os.Stat("grab_test1.go")
	t.Log(os.IsNotExist(e1))
	t.Log(os.IsExist(e1))
	t.Log(info1, e1)
}
