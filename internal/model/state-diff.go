package model

import (
	"fmt"
)

type DiffState struct {
	ReposToAdd             []FlatpakRepo
	ReposToRemove          []FlatpakRepo
	AppsToAdd              []FlatpakApplication
	AppsToRemove           []FlatpakApplication
	AppsToReplaceOverrides []FlatpakApplication
}

// StringExistsInArray checks if a string exists in an array of strings
func StringExistsInArray(target string, arr []string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}

// Function to compare Flatpak Repositories
func compareRepositories(currentRepos []FlatpakRepo, nextRepos []FlatpakRepo) ([]FlatpakRepo, []FlatpakRepo) {
	var reposToAdd []FlatpakRepo
	var reposToRemove []FlatpakRepo

	nextRepoMap := make(map[string]bool)
	for _, repo := range nextRepos {
		key := fmt.Sprintf("%s|%s|%s", repo.Name, repo.URI, repo.InstallationType)
		nextRepoMap[key] = true
	}

	for _, repo := range currentRepos {
		key := fmt.Sprintf("%s|%s|%s", repo.Name, repo.URI, repo.InstallationType)
		if _, exists := nextRepoMap[key]; !exists {
			reposToRemove = append(reposToRemove, repo)
		}
	}

	currentRepoMap := make(map[string]bool)
	for _, repo := range currentRepos {
		key := fmt.Sprintf("%s|%s|%s", repo.Name, repo.URI, repo.InstallationType)
		currentRepoMap[key] = true
	}

	for _, repo := range nextRepos {
		key := fmt.Sprintf("%s|%s|%s", repo.Name, repo.URI, repo.InstallationType)
		if _, exists := currentRepoMap[key]; !exists {
			reposToAdd = append(reposToAdd, repo)
		}
	}

	return reposToAdd, reposToRemove
}

// Function to compare Flatpak Applications
func compareApplications(currentApps []FlatpakApplication, nextApps []FlatpakApplication) ([]FlatpakApplication, []FlatpakApplication) {
	var appsToAdd []FlatpakApplication
	var appsToRemove []FlatpakApplication

	nextAppMap := make(map[string]bool)
	for _, app := range nextApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		nextAppMap[key] = true
	}

	for _, app := range currentApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		if _, exists := nextAppMap[key]; !exists {
			appsToRemove = append(appsToRemove, app)
		}
	}

	currentAppMap := make(map[string]bool)
	for _, app := range currentApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		currentAppMap[key] = true
	}

	for _, app := range nextApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		if _, exists := currentAppMap[key]; !exists {
			appsToAdd = append(appsToAdd, app)
		}
	}

	return appsToAdd, appsToRemove
}

// Function to compare Flatpak Application Permissions (Overrides)
func comparePermissions(currentApps []FlatpakApplication, nextApps []FlatpakApplication) []FlatpakApplication {
	var appsToReplace []FlatpakApplication

	for _, nextApp := range nextApps {
		for i, currentApp := range currentApps {
			if nextApp.Name == currentApp.Name && nextApp.Repo == currentApp.Repo && nextApp.InstallationType == currentApp.InstallationType {
				// Compare overrides
				overridesChanged := false
				app := FlatpakApplication{
					Name:             nextApp.Name,
					Repo:             nextApp.Repo,
					InstallationType: nextApp.InstallationType,
				}
				for _, value := range nextApp.Overrides {
					if !StringExistsInArray(value, currentApp.Overrides) {
						overridesChanged = true
						app.Overrides = append(app.Overrides, value)
						break
					}
				}

				if overridesChanged {
					appsToReplace = append(appsToReplace, app)
					// Replace overrides in the current state
					currentApps[i].Overrides = nextApp.Overrides
				}
			}
		}
	}

	return appsToReplace
}

func GetDiffState(currentState, nextState State) DiffState {
	// Handle differences as needed...
	// Compare repositories
	reposToAdd, reposToRemove := compareRepositories(currentState.Repos, nextState.Repos)
	// Compare applications
	appsToAdd, appsToRemove := compareApplications(currentState.Applications, nextState.Applications)
	// Compare permissions
	appsToReplace := comparePermissions(currentState.Applications, nextState.Applications)
	// Create FlatpakDiff structure
	return DiffState{
		ReposToAdd:             reposToAdd,
		ReposToRemove:          reposToRemove,
		AppsToAdd:              appsToAdd,
		AppsToRemove:           appsToRemove,
		AppsToReplaceOverrides: appsToReplace,
	}

}
