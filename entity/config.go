package entity

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type ActionSettings struct {
	// For readdir
	IncludeParents int `yaml:"include_parents"`

	// For readdir
	// Filter for files with given extension
	Extension string `yaml:"extension"`

	// For exec and walkdir
	// If a line starts with the given prefix, remove that prefix for building the display name.
	RemovePrefix string `yaml:"remove_prefix"`
	RemoveSuffix string `yaml:"remove_suffix"`
}

type Selector struct {
	Action string `yaml:"action"`
	// To override Id as the bucket name.
	Bucket string `yaml:"bucket"`
	// On selection execute the resulting string instead of just printing it.
	ExecOnSelect bool `yaml:"exec_on_select"`
	// After printing the selection exit with this code.
	ExitOnSelect int            `yaml:"exit_on_select"`
	Id           string         `yaml:"id"`
	OnSelect     string         `yaml:"on_select"`
	Settings     ActionSettings `yaml:"settings"`
	// Write the selection to stderr.
	StderrOutput bool     `yaml:"stderr_output"`
	Target       []string `yaml:"target"`
}

type Config struct {
	Options   Options    `yaml:"options"`
	Selectors []Selector `yaml:"selectors"`
}
