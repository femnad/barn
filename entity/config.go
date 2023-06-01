package entity

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type ActionSettings struct {
	// For readdir
	IncludeParents int `yaml:"include_parents"`

	// For exec
	// If a line starts with the users home path, remove that prefix for building the display name
	RemoveHomePrefix bool `yaml:"remove_home_prefix"`
}

type Selector struct {
	Action string `yaml:"action"`
	// On selection execute the resulting string instead of just printing it.
	ExecOnSelect bool           `yaml:"exec_on_select"`
	Id           string         `yaml:"id"`
	OnSelect     string         `yaml:"on_select"`
	Settings     ActionSettings `yaml:"settings"`
	Target       []string       `yaml:"target"`
}

type Config struct {
	Options   Options    `yaml:"options"`
	Selectors []Selector `yaml:"selectors"`
}
