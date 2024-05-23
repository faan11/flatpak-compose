package model

func GetComposeState(fileState, systemState State) State {
	var newState State

	// Iterate through fileState applications
	for _, fileApp := range fileState.Applications {
		var found bool
		var foundApp FlatpakApplication

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

	// Add repositories from fileState to newState
	for _, fileRepo := range fileState.Repos {
		var found bool
		var foundRepo FlatpakRepo
		// Search for matching repo in systemState
		for _, sysRepo := range systemState.Repos {
			if fileRepo.Name == sysRepo.Name &&
				fileRepo.Options == sysRepo.Options {
				found = true
				foundRepo = sysRepo
				break
			}
		}

		// Add the repository to newState if not found in systemState
		if found {
			newState.Repos = append(newState.Repos, foundRepo)
		}
	}

	return newState
}
