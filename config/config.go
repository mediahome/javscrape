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
	// Config ...
	// @Description: output scrape data to path
	Output string `json:"output"`
	// Config ...
	// @Description: open or close debug mode
	Debug bool `json:"debug"`
}

func DefaultConfig() *Config {
	return &Config{
		Cache:   "tmp",
		ToUpper: true,
		Output:  "output",
		Debug:   false,
	}
}
