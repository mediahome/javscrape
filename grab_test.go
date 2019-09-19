package scrape

import "testing"

// TestGrabBP4X_Find ...
func TestGrabBP4X_Find(t *testing.T) {
	grab := NewGrabBP4X(BP4XTypeJAV)
	err := grab.Find("abp-874")
	t.Log(err)
}

// TestGrabJAVBUS_Find ...
func TestGrabJAVBUS_Find(t *testing.T) {
	grab := NewGrabJAVBUS(LanguageJapanese)
	err := grab.Find("abp-874")
	t.Log(err)
}
