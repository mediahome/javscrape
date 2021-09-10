package core

import (
	"github.com/goextension/log"
)

type dummyLog struct {
}

func (d *dummyLog) Debug(args ...interface{}) {

}

func (d *dummyLog) Info(args ...interface{}) {

}

func (d *dummyLog) Warn(args ...interface{}) {

}

func (d *dummyLog) Error(args ...interface{}) {

}

func (d *dummyLog) DPanic(args ...interface{}) {

}

func (d *dummyLog) Panic(args ...interface{}) {

}

func (d *dummyLog) Fatal(args ...interface{}) {

}

func (d *dummyLog) Debugf(template string, args ...interface{}) {

}

func (d *dummyLog) Infof(template string, args ...interface{}) {

}

func (d *dummyLog) Warnf(template string, args ...interface{}) {

}

func (d *dummyLog) Errorf(template string, args ...interface{}) {

}

func (d *dummyLog) DPanicf(template string, args ...interface{}) {

}

func (d *dummyLog) Panicf(template string, args ...interface{}) {

}

func (d *dummyLog) Fatalf(template string, args ...interface{}) {

}

func (d *dummyLog) Debugw(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) Infow(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) Warnw(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) Errorw(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) DPanicw(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) Panicw(msg string, keysAndValues ...interface{}) {

}

func (d *dummyLog) Fatalw(msg string, keysAndValues ...interface{}) {

}

var NilLogger log.Logger = (*dummyLog)(nil)
