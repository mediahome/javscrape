package scrape

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const search = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	sample   bool
	details  []*javdbSearchDetail
}
