package scrape

import (
	"testing"

	"github.com/goextension/log"
	"github.com/goextension/log/zap"
	"github.com/javscrape/go-scrape/net"
)

func init() {
	zap.InitZapFileSugar()

}

// TestNewScrape ...
func TestNewScrape(t *testing.T) {
	e := net.RegisterProxy("socks5://localhost:11080")
	if e != nil {
		return
	}
	//grab1 := NewGrabBp4x(GrabBp4xTypeOption(BP4XTypeJAV))
	grab2 := NewGrabJavbus(JavbusExact(false))
	grab3 := NewGrabJavdb(JavdbExact(false))
	//doc, err := grab.Find("abp-874")
	//if err != nil {
	//	t.Fatal(err)
	scrape := NewScrape(grab2, grab3)
	scrape.Output("video")
	scrape.GrabSample(true)
	scrape.ImageCache("")
	msg, e := scrape.Find("abp-890")
	if e != nil {
		t.Fatal(e)
	}
	for _, m := range *msg {
		log.Infof("%+v", m)
	}
}
