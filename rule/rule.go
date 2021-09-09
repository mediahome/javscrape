package rule

// Rule ...
// @Description:
type Rule struct {
	Entrance string   `toml:"entrance"`
	MainPage string   `toml:"main_page"`
	Actions  []Action `toml:"actions"`
}
