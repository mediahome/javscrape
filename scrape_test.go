package scrape

import (
	"log"
	"testing"
)

func init() {
	//zap.InitZapSugar()
	DebugOn()
}

// TestNewScrape ...
func TestNewScrape(t *testing.T) {

	e := RegisterProxy("socks5://localhost:10808")
	if e != nil {
		return
	}
	//grab1 := NewGrabBp4x(GrabBp4xTypeOption(BP4XTypeJAV))
	grab2 := NewGrabJavbus()
	grab3 := NewGrabJavdb()
	//doc, err := grab.Find("abp-874")
	//if err != nil {
	//	t.Fatal(err)
	scrape := NewScrape(GrabOption(grab2), GrabOption(grab3), OptimizeOption(true))
	//scrape.Output("video")
	//scrape.GrabSample(true)
	e = scrape.Find("abp-891")
	checkErr(e)
	e = scrape.Find("abp-892")
	checkErr(e)
	e = scrape.Range(func(key string, content Content) error {
		t.Log("key", key, "info", content)
		return nil
	})
	checkErr(e)
	e = scrape.Output()
	checkErr(e)
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
