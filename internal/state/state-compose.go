package state 

import "github.com/faan11/flatpak-compose/internal/model"


func compareRemoteEnvironments(env1, env2 model.Environment) *model.Environment {
	if env1.InstallationType != env2.InstallationType {
		return nil
	}

	/*commonCore := make(map[string]string)
	for k, v := range env1.Core {
		if v2, exists := env2.Core[k]; exists && v == v2 {
			commonCore[k] = v
		}
	}*/

	commonRemotes := make(map[string]map[string]string)
	for k, _ := range env1.Remotes {
		if v2, exists := env2.Remotes[k]; exists  {
			commonRemotes[k] = v2
		}
	}

	if len(commonRemotes) == 0 {
		return nil
	}

	return &model.Environment{
		Core:            make(map[string]string),
		Remotes:         commonRemotes,
		InstallationType: env1.InstallationType,
	}
}

func GetComposeState(fileState, systemState model.State) model.State {
	var newState model.State

	// Iterate through fileState applications
	for _, fileApp := range fileState.Applications {
		var found bool
		var foundApp model.FlatpakApplication

		// Search for matching app in systemState
		for _, sysApp := range systemState.Applications {
			if fileApp.Name == sysApp.Name &&
				fileApp.Repo == sysApp.Repo &&
				fileApp.Branch == sysApp.Branch &&
				fileApp.InstallationType == sysApp.InstallationType {
				found = true
				foundApp = sysApp
				break
			}
		}

		// Add the application to the newState if not found in systemState
		if found {
			newState.Applications = append(newState.Applications, foundApp)
		}
	}

	for _, fileEnv := range fileState.Environment {
		for _, sysEnv := range systemState.Environment {
			env := compareRemoteEnvironments(fileEnv,sysEnv)
			// Add repositories from fileState to newState
			if env != nil {
				newState.Environment = append(newState.Environment, *env)
				break;
			}
		}
	}

	/*for _, fileRepo := range fileState.Environment {
		var found bool
		var foundRepo model.Environment
		// Search for matching repo in systemState
		for _, sysRepo := range systemState.Environment {
			if fileRepo.Name == sysRepo.Name &&
				fileRepo.InstallationType == sysRepo.InstallationType {
				found = true
				foundRepo = sysRepo
				break
			}
		}

		// Add the repository to newState if not found in systemState
		if found {
			newState.Repos = append(newState.Repos, foundRepo)
		}
	}*/

	return newState
}
