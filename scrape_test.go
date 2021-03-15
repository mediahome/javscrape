package scrape

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
)

func init() {
	//zap.InitZapSugar()
	DefaultOutputPath = `D:\workspace\golang\project\go-scrape\video`
	//DebugOn()
	debug = true
}

// TestNewScrape ...
func TestNewScrape(t *testing.T) {
	var e error
	DebugOn()
	e = RegisterProxy("http://localhost:7890")
	if e != nil {
		return
	}
	//grab1 := NewGrabBp4x(GrabBp4xTypeOption(BP4XTypeJAV))
	grab2 := NewGrabJavbus()
	//grab2.SetLanguage(LanguageEnglish)
	grab3 := NewGrabJavdb()
	//doc, err := grab.Find("abp-874")
	//if err != nil {
	//	t.Fatal(err)
	scrape := NewScrape(GrabOption(grab2), GrabOption(grab3), ExactOption(true))
	//scrape.Output("video")
	//scrape.GrabSample(true)
	e = scrape.Find("HMDN-344")
	checkErr(e)
	scrape.Range(func(key string, content Content) error {
		fmt.Printf("key:%v,content:%+v", key, content)
		return nil
	})
	e = scrape.Output()
	checkErr(e)
	scrape.Clear()
	//e = scrape.Find("snis")
	//checkErr(e)
	//e = scrape.Output()
	//checkErr(e)
	//scrape.Clear()
	//e = scrape.Find("ssni")
	//checkErr(e)
	//e = scrape.Output()
	//checkErr(e)
	//scrape.Clear()
}

// TestNewScrape ...
func TestNewScrapeOutput(t *testing.T) {
	var e error
	DebugOn()
	e = RegisterProxy("http://localhost:7890")
	if e != nil {
		return
	}
	//grab1 := NewGrabBp4x(GrabBp4xTypeOption(BP4XTypeJAV))
	grab2 := NewGrabJavbus()
	//grab2.SetLanguage(LanguageEnglish)
	grab3 := NewGrabJavdb()
	//doc, err := grab.Find("abp-874")
	//if err != nil {
	//	t.Fatal(err)
	scrape := NewScrape(GrabOption(grab2), GrabOption(grab3), ExactOption(false))
	//scrape.Output("video")
	//scrape.GrabSample(true)
	e = scrape.Find("HMDN-344")
	checkErr(e)
	scrape.Range(func(key string, content Content) error {
		log.Printf("key:%v,content:%+v", key, content)
		return nil
	})
	outputFlag := "javdb"

	scrape.OutputCallback(func(key string, content Content) *OutputOption {
		option := DefaultOutputOption()
		option.OutputPath = filepath.Join(DefaultOutputPath, key)
		option.CopyInfo = true
		option.InfoName = key + ".nfo"
		if outputFlag != content.From {
			option.Skip = true
		}
		return option
	})
	checkErr(e)
	scrape.Clear()
	//e = scrape.Find("snis")
	//checkErr(e)
	//e = scrape.Output()
	//checkErr(e)
	//scrape.Clear()
	//e = scrape.Find("ssni")
	//checkErr(e)
	//e = scrape.Output()
	//checkErr(e)
	//scrape.Clear()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
