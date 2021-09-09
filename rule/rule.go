package rule

type Rule struct {
	Entrance string `toml:"entrance"`
	MainPage string `toml:"main_page"`
	Actions  Action `toml:"actions"`
	//ActionGroup map[string][]Action `toml:"action_group"`
}
