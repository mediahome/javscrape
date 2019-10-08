package scrape

import "testing"

// TestNewJavdb ...
func TestNewJavdb(t *testing.T) {
	DebugOn()
	javdb := NewJavdb()
	grab, e := javdb.Find("heyzo-2097")
	t.Log(grab, e)
}
