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

func NewStringValue(value interface{}) *Value {
	return &Value{Type: rule.ProcessValueString, v: value}
}

func NewArrayValue(value []interface{}) *Value {
	return &Value{Type: rule.ProcessValueString, v: value}
}

func NewMapValue(value interface{}) *Value {
	return &Value{Type: rule.ProcessValueMap, v: value}
}

func NewFileValue(value []byte, fn func(key string, value []byte)) *Value {
	key := MD5(value)
	fn(key, value)
	return &Value{Type: rule.ProcessValueFie, v: value}
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
