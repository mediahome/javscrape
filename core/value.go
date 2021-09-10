package core

import (
	"fmt"

	"github.com/goextension/gomap"

	"github.com/javscrape/go-scrape/rule"
)

type Value struct {
	Type rule.ProcessValue
	v    interface{}
}

func (v *Value) Set(value interface{}) {
	v.v = value
}

func (v *Value) SetFile(value []byte, fn func(key string, value []byte)) {
	key := MD5(value)
	v.v = key
	fn(key, value)
}

func (v Value) GetMap() gomap.Map {
	return v.v.(gomap.Map)
}

func (v Value) GetArray() []string {
	return v.v.([]string)
}

func (v Value) GetString() string {
	return v.v.(string)
}

func (v Value) GetFileHash() string {
	return v.v.(string)
}

func (v Value) String() string {
	return fmt.Sprintf("Value(Type:%v,Value:%+v)", v.Type, v.v)
}
