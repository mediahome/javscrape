package scrape

import "testing"

// TestGrabBP4X_Find ...
func TestGrabBP4X_Find(t *testing.T) {
	grab := NewGrabBP4X(BP4XTypeJAV)
	err := grab.Find("abp-874")
	t.Log(err)
}
