package rule

type Value struct {
	Name  string      `toml:"name"`
	Value interface{} `toml:"value"`
}
