package config

type Selector struct {
	Action   string `yaml:"action"`
	Arg      string `yaml:"arg"`
	Id       string `yaml:"id"`
	OnSelect string `yaml:"on_select"`
}

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type Config struct {
	Selectors []Selector `yaml:"selectors"`
	Options   Options    `yaml:"options"`
}
