package rule

//Action ...
type Action struct {
	Type      ActionType `toml:"type,omitempty"`
	Name      string     `toml:"name,omitempty"`
	Index     int        `toml:"index,omitempty"`
	Web       Web        `toml:"web,omitempty"`
	Through   bool       `toml:"through,omitempty"`
	OnSuccess string     `toml:"on_success,omitempty"`
	OnFailure string     `toml:"on_failure,omitempty"`
	Success   []Process  `toml:"success,omitempty"`
}
