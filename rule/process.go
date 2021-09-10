package rule

type Process struct {
	Name     string      `toml:"name"`
	Trim     bool        `toml:"trim"`
	Type     ProcessType `toml:"type"`
	Property string      `toml:"property"`
}
