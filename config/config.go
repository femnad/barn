package config

type Selector struct {
	Action   string   `yaml:"action"`
	Args     []string `yaml:"args"`
	Id       string   `yaml:"id"`
	OnSelect string   `yaml:"on_select"`
}

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type Config struct {
	Selectors []Selector `yaml:"selectors"`
	Options   Options    `yaml:"options"`
}
