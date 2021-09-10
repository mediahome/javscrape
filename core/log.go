package core

import (
	"github.com/goextension/log/zap"
)

func InitGlobalLogger(debug bool) {
	DEBUG = debug
	if debug {
		zap.InitZapSugar()
		return
	}
}
