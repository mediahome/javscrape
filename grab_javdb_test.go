package scrape

import "testing"

// TestNewJavdb ...
func TestNewJavdb(t *testing.T) {
	DebugOn()
	javdb := NewJavdb()
	grab, e := javdb.Find("abp")
	t.Log(grab, e)
	for grab.HasNext() {
		iGrab, e := grab.Next()
		t.Log(iGrab, e)
	}

}
