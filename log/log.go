package log

import (
	"github.com/goextension/log"
	"github.com/goextension/log/zap"

	"github.com/javscrape/go-scrape/core"
)

func InitGlobalLogger(debug bool) {
	core.DEBUG = debug
	if debug {
		zap.InitZapSugar()

		return
	}
	log.Register(&dummy{})
}
