package core

import (
	"fmt"
	"os"

	"github.com/goextension/log"
)

type printLog struct {
}

func (p *printLog) Debug(args ...interface{}) {
	fmt.Print("[Debug] ")
	fmt.Println(args...)
}

func (p *printLog) Info(args ...interface{}) {
	fmt.Print("[Info] ")
	fmt.Println(args...)
}

func (p *printLog) Warn(args ...interface{}) {
	fmt.Print("[Warn] ")
	fmt.Println(args...)
}

func (p *printLog) Error(args ...interface{}) {
	fmt.Print("[Error] ")
	fmt.Println(args...)
}

func (p *printLog) DPanic(args ...interface{}) {
	fmt.Print("[DPanic] ")
	fmt.Println(args...)
	panic("dpanic call")
}

func (p *printLog) Panic(args ...interface{}) {
	fmt.Print("[Panic] ")
	fmt.Println(args...)
	panic("panic call")
}

func (p *printLog) Fatal(args ...interface{}) {
	fmt.Print("[Fatal] ")
	fmt.Println(args...)
	os.Exit(0)
}

func (p *printLog) Debugf(template string, args ...interface{}) {
	fmt.Print("[Debugf] ")
	fmt.Printf(template+"\n", args...)
}

func (p *printLog) Infof(template string, args ...interface{}) {
	fmt.Print("[Infof] ")
	fmt.Printf(template+"\n", args...)
}

func (p *printLog) Warnf(template string, args ...interface{}) {
	fmt.Print("[Warnf] ")
	fmt.Printf(template+"\n", args...)
}

func (p *printLog) Errorf(template string, args ...interface{}) {
	fmt.Print("[Errorf] ")
	fmt.Printf(template+"\n", args...)
}

func (p *printLog) DPanicf(template string, args ...interface{}) {
	fmt.Print("[DPanicf] ")
	fmt.Printf(template+"\n", args...)
	panic("dpanicf call")
}

func (p *printLog) Panicf(template string, args ...interface{}) {
	fmt.Print("[Panicf] ")
	fmt.Printf(template+"\n", args...)
	panic("panicf call")
}

func (p *printLog) Fatalf(template string, args ...interface{}) {
	fmt.Print("[Fatalf] ")
	fmt.Printf(template+"\n", args...)
	os.Exit(0)
}

func (p *printLog) Debugw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Debugw] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Infow] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) Warnw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Warnw] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) Errorw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Errorw] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) DPanicw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[DPanicw] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) Panicw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Panicw] ", msg)
	fmt.Println(keysAndValues...)
}

func (p *printLog) Fatalw(msg string, keysAndValues ...interface{}) {
	fmt.Print("[Fatalw] ", msg)
	fmt.Println(keysAndValues...)
}

var PrintLogger log.Logger = (*printLog)(nil)
