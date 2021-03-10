package scrape

import "testing"

// TestNewGrabJAVBUS ...
func TestNewGrabJAVBUS(t *testing.T) {
	//DebugOn()
	javbus := NewGrabJavbus()
	javbus.SetSample(true)
	//javbus.SetLanguage(LanguageEnglish)
	grab, e := javbus.Find("vec-457")
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
