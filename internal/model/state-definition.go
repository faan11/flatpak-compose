package model

// Environment has core + remotes of an installation type
type Environment struct {
	Core    map[string]string		`yaml:"core"`
	Remotes map[string]map[string]string	`yaml:"remotes"`
	InstallationType string 		`yaml:"type"`
}

type Permission struct {
	Table		string		`yaml:"table"`		 
	Object		string		`yaml:"object"`
	Permission	string		`yaml:"permission"`
	Data		string		`yaml:"data"`
}

type FlatpakApplication struct {
	Name             string   	`yaml:"name"`  
	Repo             string   	`yaml:"repo"`
	Branch           string   	`yaml:"branch,omitempty"`
	All              []string 	`yaml:"all"`            // Default permissions
	Overrides        []string 	`yaml:"overrides"`      // Override permissions
	OverridesUser    []string 	`yaml:"overrides_user"` // Override user permissions
	InstallationType string   	`yaml:"type"`
	Permissions	 []Permission	`yaml:"permissions"`
}

type State struct {
	Environment  []Environment 	  `yaml:"envs"`
	Applications []FlatpakApplication `yaml:"applications"`
}


func (e Environment) RemoteExists(installationType string, remoteName string) bool {
	if (installationType == e.InstallationType){
		_, exists := e.Remotes[remoteName]
		return exists
	} else {
		return false;
	}
}
