package log

type dummy struct {
}

func (d dummy) Debug(args ...interface{}) {

}

func (d dummy) Info(args ...interface{}) {

}

func (d dummy) Warn(args ...interface{}) {

}

func (d dummy) Error(args ...interface{}) {

}

func (d dummy) DPanic(args ...interface{}) {

}

func (d dummy) Panic(args ...interface{}) {

}

func (d dummy) Fatal(args ...interface{}) {

}

func (d dummy) Debugf(template string, args ...interface{}) {

}

func (d dummy) Infof(template string, args ...interface{}) {

}

func (d dummy) Warnf(template string, args ...interface{}) {

}

func (d dummy) Errorf(template string, args ...interface{}) {

}

func (d dummy) DPanicf(template string, args ...interface{}) {

}

func (d dummy) Panicf(template string, args ...interface{}) {

}

func (d dummy) Fatalf(template string, args ...interface{}) {

}

func (d dummy) Debugw(msg string, keysAndValues ...interface{}) {

}

func (d dummy) Infow(msg string, keysAndValues ...interface{}) {

}

func (d dummy) Warnw(msg string, keysAndValues ...interface{}) {

}

func (d dummy) Errorw(msg string, keysAndValues ...interface{}) {

}

func (d dummy) DPanicw(msg string, keysAndValues ...interface{}) {

}

func (d dummy) Panicw(msg string, keysAndValues ...interface{}) {

}

func (d dummy) Fatalw(msg string, keysAndValues ...interface{}) {

}
