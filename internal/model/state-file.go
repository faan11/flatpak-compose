package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func GetFileState(stateFile string) (State, error) {
	var config State

	yamlFile, err := os.ReadFile(stateFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	// Check for at least one repo
	if len(config.Repos) == 0 {
		return config, fmt.Errorf("at least one repository must exist")
	}

	for i := range config.Repos {
		if config.Repos[i].Options == "" {
			config.Repos[i].Options = "system" // or any default value you prefer
		}
	}

	repoNames := make(map[string]bool)
	for _, repo := range config.Repos {
		// Check for unique repo names
		key := repo.Name + "|" + repo.Options
		if repoNames[key] {
			return config, fmt.Errorf("duplicate repository name found: %s", key)
		}
		repoNames[key] = true
	}

	for i, app := range config.Applications {
		// Set the default repo if not specified and only one repo exists
		if len(config.Repos) == 1 && app.Repo == "" {
			config.Applications[i].Repo = config.Repos[0].Name
		}

		// Set default branch to "stable" if not specified
		if app.Branch == "" {
			config.Applications[i].Branch = "stable"
		}

		// Check if the repo specified for an application exists in the list of repos
		repoExists := false
		for _, repo := range config.Repos {
			if app.Repo == repo.Name && app.InstallationType == repo.Options {
				repoExists = true
				break
			}
		}
		if !repoExists {
			return config, fmt.Errorf("Application '%s' refers to a non-existent repository: %s in %s mode", app.Name, app.Repo, app.InstallationType)
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
