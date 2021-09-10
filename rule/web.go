package rule

type Web struct {
	Method   string              `toml:"method,omitempty"`
	Header   map[string][]string `toml:"header,omitempty"`
	URL      string              `toml:"url,omitempty"`
	URI      string              `toml:"uri,omitempty"`
	Selector string              `toml:"selector,omitempty"`
	Success  []Process           `toml:"success,omitempty"`
}
