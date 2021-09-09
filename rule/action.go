package rule

//Action ...
type Action struct {
	Name      string     `toml:"name"`
	Index     int        `toml:"index"`
	Type      ActionType `toml:"type"`
	URI       string     `toml:"uri"`
	Through   bool       `toml:"through"`
	OnSuccess string     `toml:"on_success"`
	OnFailure string     `toml:"on_failure"`
}
