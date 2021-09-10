package scrape

import (
	"testing"

	"github.com/javscrape/go-scrape/config"
	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/rule"
)

var cfg = config.DefaultConfig()
var scrape core.IScrape

func init() {
	cfg.Debug = true
	scrape = New(ProxyOption("http://localhost:7890"))
}

// TestNewScrape ...
func TestNew(t *testing.T) {
	r, err := rule.LoadRuleFromFile("tmp.toml")
	if err != nil {
		t.Fatal(err)
	}

	grab, err := scrape.LoadRules(r)
	if err != nil {
		t.Fatal(err)
		return
	}

	grab[0].Do()
}
