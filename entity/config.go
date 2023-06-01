package entity

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type ActionArgs struct {
	// For readdir
	IncludeParents int `yaml:"include_parents"`
}

type Selector struct {
	Action   string     `yaml:"action"`
	Args     ActionArgs `yaml:"args"`
	Id       string     `yaml:"id"`
	OnSelect string     `yaml:"on_select"`
	Path     []string   `yaml:"path"`
}

type Config struct {
	Options   Options    `yaml:"options"`
	Selectors []Selector `yaml:"selectors"`
}
