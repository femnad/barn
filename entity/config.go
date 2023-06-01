package entity

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type Selector struct {
	Action   string   `yaml:"action"`
	Args     []string `yaml:"args"`
	Id       string   `yaml:"id"`
	OnSelect string   `yaml:"on_select"`
}

type Config struct {
	Options   Options    `yaml:"options"`
	Selectors []Selector `yaml:"selectors"`
}
