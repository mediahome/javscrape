package rule

type Web struct {
	Method   string              `toml:"method,omitempty"`
	Header   map[string][]string `toml:"header,omitempty"`
	Value    []string            `toml:"value,omitempty"`
	Relative bool                `toml:"relative,omitempty"`
	Skip     []SkipType          `toml:"skip,omitempty"`
	Selector string              `toml:"selector,omitempty"`
	Success  []Process           `toml:"success,omitempty"`
}
