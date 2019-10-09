package scrape

import (
	"testing"
)

// TestNewScrape ...
func TestNewScrape(t *testing.T) {
	//e := net.RegisterProxy("socks5://localhost:11080")
	//if e != nil {
	//	return
	//}
	//grab1 := NewGrabBp4x(GrabBp4xTypeOption(BP4XTypeJAV))
	grab2 := NewGrabJavbus()
	grab3 := NewGrabJavdb()
	//doc, err := grab.Find("abp-874")
	//if err != nil {
	//	t.Fatal(err)
	scrape := NewScrape(grab2, grab3)

	scrape.GrabSample(true)

	msg, e := scrape.Find("abp-874")
	if e != nil {
		return
	}
	for _, m := range *msg {
		log.Infof("%+v", m)
	}
}
