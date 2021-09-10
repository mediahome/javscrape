package rule

//Action ...
type Action struct {
	Type  ActionType `toml:"type"`
	Name  string     `toml:"name"`
	Index int        `toml:"index"`
	Web   Web        `toml:"web"`

	Step      Step   `toml:"step"`
	Through   bool   `toml:"through"`
	OnSuccess string `toml:"on_success"`
	OnFailure string `toml:"on_failure"`
}
