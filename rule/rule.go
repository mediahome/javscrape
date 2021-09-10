package rule

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Rule ...
// @Description:
type Rule struct {
	Entrance  string    `toml:"entrance,omitempty"`
	MainPage  string    `toml:"main_page,omitempty"`
	InputKey  string    `toml:"input_key,omitempty"`
	InputType InputType `toml:"input_type,omitempty"`
	Actions   []Action  `toml:"actions,omitempty"`
}

func LoadRuleFromFile(file string) (*Rule, error) {
	var r Rule
	_, err := toml.DecodeFile(file, &r)
	return &r, err
}

func SaveRuleToFile(file string, r *Rule) error {
	openFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_SYNC, 0755)
	if err != nil {
		return err
	}
	defer openFile.Close()
	enc := toml.NewEncoder(openFile)
	return enc.Encode(r)
}
