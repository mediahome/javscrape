package scrape

import (
	"log"
	"testing"
)

func init() {
	//zap.InitZapSugar()
	DefaultOutputPath = `D:\workspace\golang\project\go-scrape\video`
	//DebugOn()
	debug = false
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
	scrape := NewScrape(GrabOption(grab2), GrabOption(grab3), OptimizeOption(true), ExactOption(false))
	//scrape.Output("video")
	//scrape.GrabSample(true)
	e = scrape.Find("abp")
	checkErr(e)
	e = scrape.Find("snis")
	checkErr(e)
	e = scrape.Find("ssni")
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
