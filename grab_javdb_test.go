package scrape

import "testing"

// TestNewJavdb ...
func TestNewJavdb(t *testing.T) {
	DebugOn()
	javdb := NewGrabJavdb()
	javdb.Sample(true)
	grab, e := javdb.Find("FCH-041")
	t.Log(grab, e)
	t.Log(grab.HasNext())
	count := 0
	for grab.HasNext() {
		if count > 2 {
			return
		}
		grab, e = grab.Next()
		t.Log(grab, e)
		count++
	}

}
