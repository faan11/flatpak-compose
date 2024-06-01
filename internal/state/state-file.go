package state 

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"github.com/faan11/flatpak-compose/internal/model"
)

func GetFileState(stateFile string) (model.State, error) {
	var config model.State

	yamlFile, err := os.ReadFile(stateFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	for i,_ := range config.Environment {	
		if config.Environment[i].InstallationType == "" {
			config.Environment[i].InstallationType = "system" // or any default value you prefer
		}
	}

	repoNames := make(map[string]bool)
	for _, env := range config.Environment {
		for name, _ := range env.Remotes {
			// Check for unique repo names
			key := name + "|" + env.InstallationType
			if repoNames[key] {
				return config, fmt.Errorf("duplicate repository name found: %s", key)
			}
			repoNames[key] = true
		}
	}

	for i, app := range config.Applications {
		// Set the default repo if not specified and only one repo exists
		if len(config.Environment) >= 1 && len(config.Environment[0].Remotes) >= 1 && app.Repo == "" {
			rem := config.Environment[0].Remotes
			// Collect all keys into a slice
	       	        keys := make([]string, 0, len(rem))
			for k,_ := range rem {
				keys = append(keys, k)
			}
		    	// Get a random key from the slice
		    	randomKey := keys[0]
			config.Applications[i].Repo = randomKey 
		}

		// Set default branch to "stable" if not specified
		if app.Branch == "" {
			config.Applications[i].Branch = "stable"
		}

		// Check if the repo specified for an application exists in the list of repos
		repoExists := false
		appRepo := app.Repo + "|" + app.InstallationType
		for k, v := range repoNames {
			if v == true && k == appRepo {
				repoExists = true
				break
			}
		}
		if !repoExists {
			fmt.Printf("Warning: application '%s' refers to a non-existent repository: '%s' in '%s' mode and will be ignored during installation process but overrides will still be applied if possible\n", app.Name, app.Repo, app.InstallationType)
		}

		// Check for valid InstallationType
		if app.InstallationType != "user" && app.InstallationType != "system" {
			return config, fmt.Errorf("application '%s' has an invalid InstallationType: %s", app.Name, app.InstallationType)
		}

		// Ensure application name is required
		if app.Name == "" {
			return config, fmt.Errorf("application name is required")
		}
	}
	return config, nil
}
