package scrape

// DefaultJavdbMainPage ...
const DefaultJavdbMainPage = "https://javdb2.com"
const javdbSearch = "/search?q=%s&f=all"

type grabJavdb struct {
	mainPage string
	sample   bool
	//details  []*javdbSearchDetail
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
func (g *grabJavdb) Find(name string) (IGrab, error) {
	url := g.mainPage + javdbSearch
	results, e := javdbSearchResultAnalyze(url, name)
	if e != nil {
		return nil, e
	}
	log.Info(results)
	return g, nil
}

type javdbSearchResult struct {
}

func javdbSearchResultAnalyze(url, name string) (result *javdbSearchResult, e error) {
	return &javdbSearchResult{}, nil
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
	}
}
