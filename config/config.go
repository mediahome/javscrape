package config

// Config ...
// @Description: scrape configuration
type Config struct {
	// Config ...
	// @Description: cache path
	Cache string `json:"cache"`
	// Config ...
	// @Description:case id to upper
	ToUpper bool `json:"to_upper"`
}

func DefaultConfig() *Config {
	return &Config{
		Cache:   "tmp",
		ToUpper: true,
	}
}
