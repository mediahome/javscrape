package scrape

import (
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/core"
	"github.com/javscrape/go-scrape/internal"
)

func init() {
	log.Register(core.NiLogger)
}

var New = internal.NewScrape
