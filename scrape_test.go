package scrape

import (
	"testing"

	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/config"
	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

var cfg = config.DefaultConfig()
var scrape core.IScrape

func init() {
	cfg.Debug = true
	core.DEBUG = true
	//core.InitGlobalLogger(cfg.Debug)
}

// TestNewScrape ...
func TestNew(t *testing.T) {
	scrape = New(ProxyOption("http://127.0.0.1:7890"))

	r, err := rule.LoadRuleFromFile("./templates/javbus.toml")
	if err != nil {
		t.Fatal(err)
	}
	log.Debug("TEST", "load rules")
	grabs, err := scrape.LoadRules(r)
	if err != nil {
		t.Fatal(err)
	}

	if len(grabs) == 0 {
		t.Fatal("empty grabs list")
	}

	err = grabs[0].Do("ABW-140")
	if err != nil {
		t.Fatal(err)
		return
	}

	log.Debug("Test", "total values", grabs[0].Value())
}
