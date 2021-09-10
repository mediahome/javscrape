package rule

type Web struct {
	Method   string              `toml:"method"`
	Header   map[string][]string `toml:"header"`
	URL      string              `toml:"url"`
	URI      string              `toml:"uri"`
	Selector string              `toml:"selector"`
}
