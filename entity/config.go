package entity

type Options struct {
	DatabasePath string `yaml:"db_path"`
}

type ActionSettings struct {
	// Common to all actions
	// Persist all the items in the action output.
	Eager bool `yaml:"eager"`
	// On selection execute the resulting string instead of just printing it.
	ExecOnSelect bool `yaml:"exec_on_select"`
	// After printing the selection exit with this code.
	ExitOnSelect int `yaml:"exit_on_select"`
	// Template for outputting result based on the selected entry.
	OnSelect string `yaml:"on_select"`
	// Write the selection to stderr.
	StderrOutput bool `yaml:"stderr_output"`

	// For exec
	// Command for determining the base directory before executing the action command.
	SetPwdCmd string `yaml:"set_pwd_cmd"`

	// For exec and walkdir
	// If a line starts with the given prefix, remove that prefix for building the display name.
	RemovePrefix string `yaml:"remove_prefix"`
	RemoveSuffix string `yaml:"remove_suffix"`

	// For readdir
	// Filter for files with given extension.
	Extension string `yaml:"extension"`
	// Include given number of parent dirs when persisting the selection.
	IncludeParents int `yaml:"include_parents"`
}

type Selector struct {
	Action string `yaml:"action"`
	Id     string `yaml:"id"`
	// To override Id as the bucket name.
	Bucket   string         `yaml:"bucket"`
	Settings ActionSettings `yaml:"settings"`
	Target   []string       `yaml:"target"`
}

type Config struct {
	Options   Options    `yaml:"options"`
	Selectors []Selector `yaml:"selectors"`
}
