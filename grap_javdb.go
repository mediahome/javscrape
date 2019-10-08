package scrape

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const search = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	sample   bool
	details  []*javdbSearchDetail
}

// Sample ...
func (g *grabJavdb) Sample(b bool) {
	g.sample = b
}

// Name ...
func (g *grabJavdb) Name() string {
	return "javdb"
}

// Find ...
func (g *grabJavdb) Find(string) (IGrab, error) {
	panic("implement me")
}

// Decode ...
func (g *grabJavdb) Decode(*[]*Message) error {
	panic("implement me")
}

// MainPage ...
func (g *grabJavdb) MainPage(url string) {
	g.mainPage = url
}

// NewJavdb ...
func NewJavdb() IGrab {
	return &grabJavdb{
		mainPage: DefaultJavdbMainPage,
		sample:   false,
		details:  nil,
	}
}
