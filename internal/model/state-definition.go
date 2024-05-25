package model

type FlatpakRepo struct {
	Name             string `yaml:"name"`
	URI              string `yaml:"uri"`
	InstallationType string `yaml:"type"`
}

type FlatpakApplication struct {
	Name             string   `yaml:"name"`
	Repo             string   `yaml:"repo"`
	Branch           string   `yaml:"branch,omitempty"`
	All              []string `yaml:"all"`            // Default permissions
	Overrides        []string `yaml:"overrides"`      // Override permissions
	OverridesUser    []string `yaml:"overrides_user"` // Override user permissions
	InstallationType string   `yaml:"type"`
}

type State struct {
	Repos        []FlatpakRepo        `yaml:"repos"`
	Applications []FlatpakApplication `yaml:"applications"`
}
