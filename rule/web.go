package rule

type Web struct {
	Method    string              `toml:"method,omitempty"`
	Header    map[string][]string `toml:"header,omitempty"`
	BeforeURL string              `toml:"before_url,omitempty"`
	FromValue []string            `toml:"from_value,omitempty"`
	AfterURL  string              `toml:"after_url,omitempty"`
	Relative  bool                `toml:"relative,omitempty"`
	Skip      []SkipType          `toml:"skip,omitempty"`
	Selector  string              `toml:"selector,omitempty"`
	Success   []Process           `toml:"success,omitempty"`
}
