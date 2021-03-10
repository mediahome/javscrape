package scrape

import (
	"os"
	"testing"
)

// TestGrabJAVBUS_Find ...
func TestGrabJAVBUS_Find(t *testing.T) {
	DebugOn()
	e := RegisterProxy("http://localhost:7890")
	if e != nil {
		return
	}
	grab := NewGrabJavbus()
	grab.SetSample(true)
	doc, err := grab.Find("vec-457")

	if err != nil {
		t.Fatal(err)
	}
	msg, err := doc.Result()
	t.Logf("%+v", msg)
	if err != nil {
		t.Fatal(err)
	}
	cache := newCache()

	err = imageCache(cache, msg[0], true)
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
