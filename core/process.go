package core

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/goextension/log"

	"github.com/javscrape/go-scrape/rule"
)

func ProcessValue(selection *goquery.Selection, p rule.Process) *Value {
	var v string
	switch p.Property {
	case rule.ProcessPropertyArray:
		var arr []interface{}
		selection.Each(func(i int, selection *goquery.Selection) {
			v = strings.TrimSpace(selection.Text())
			log.Debug("ACTION", "array", v)
			arr = append(arr, v)
		})
		if len(arr) != 0 {
			return NewArrayValue(arr)
		}
	case rule.ProcessPropertyValue:
		v = selection.Text()
	case rule.ProcessPropertyAttr:
		v = selection.AttrOr(p.PropertyName, "")
	case rule.ProcessPropertyText:
		selection.Contents().Each(func(i int, selection *goquery.Selection) {
			if goquery.NodeName(selection) == "#text" {
				v = selection.Text()
			}
		})
	}

	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return NewStringValue(v)
}
