package scrape

import "testing"

// TestNewJavdb ...
func TestNewJavdb(t *testing.T) {
	DebugOn()
	javdb := NewJavdb()
	grab, e := javdb.Find("abp-874")
	t.Log(grab, e)
}
