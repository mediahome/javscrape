package rule

type Process struct {
	Name         string       `toml:"name,omitempty"`
	Selector     string       `toml:"selector,omitempty"`
	Filter       []string     `toml:"filter,omitempty"`
	Index        int          `toml:"index,omitempty"`
	Trim         bool         `toml:"trim,omitempty"`
	Type         ProcessType  `toml:"type,omitempty"`
	Property     string       `toml:"property,omitempty"`
	PropertyName string       `toml:"property_name,omitempty"`
	Value        ProcessValue `toml:"value,omitempty"`
}
