package state

import (
	"fmt"
	"slices"
	"github.com/faan11/flatpak-compose/internal/model"
)

type DiffState struct {
	EnvToAdd               []model.Environment
	EnvToRemove            []model.Environment
	EnvToUpdate            []model.Environment
	AppsToAdd              []model.FlatpakApplication
	AppsToRemove           []model.FlatpakApplication
	AppsToReplaceOverrides []model.FlatpakApplication
}

type diffCore struct {
	added   map[string]string
	removed map[string]string
	updated map[string]string
}

type diffRemote struct {
	added   map[string]map[string]string
	removed map[string]map[string]string
	updated map[string]map[string]string
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

//
// Environment has core + remotes of an installation type
/*

Compares a slice of environment "prev" with a slice of environment "next".
The result of this function are three slices of environment "tobeadd","toberemoved" and "tobeupdated".

I can compare only Environment with the same InstallationType.
if prev has user environment and next does not, the environment needs to be added ENTIRELY.
if prev has no user environment and next has it, the environment needs to be removed ENTIRELY but this is an error and should not occur.

For each environment prev, next with the same installation type:

In the Core,  i need to check:
- keys that are present in "prev" but not in "next" this needs to be removed
- keys that are present in "next" but not in "prev" this needs to be added
- keys that are present in both but their value are different, needs to be updated

In the Environment,  i need to check:
- keys that are present in "prev" but not in "next" this needs to be removed
- keys that are present in "next" but not in "prev" this needs to be added
- keys that are present in both but their value are different, needs to be updated
*/


func compareEnvironments(prev, next []model.Environment) (toBeAdded, toBeRemoved, toBeUpdated []model.Environment) {
	prevMap := make(map[string]model.Environment)
	nextMap := make(map[string]model.Environment)

	// Convert slices to maps for easier comparison
	for _, env := range prev {
		prevMap[env.InstallationType] = env
	}
	for _, env := range next {
		nextMap[env.InstallationType] = env
	}

	// Compare environments
	for key, prevEnv := range prevMap {
		nextEnv, exists := nextMap[key]
		if !exists {
			// If prev has an environment that next does not have, it should be added to next
			toBeAdded = append(toBeAdded, prevEnv)
			continue
		}

		// Check core differences
		coreDiff := compareCore(prevEnv.Core, nextEnv.Core)
		remotesDiff := compareRemotes(prevEnv.Remotes, nextEnv.Remotes)

		// If there are differences, add to toBeUpdated
		
		addedEnv := model.Environment{
			Core:            coreDiff.added,
			Remotes:         remotesDiff.added,
			InstallationType: nextEnv.InstallationType,
		}

		removedEnv := model.Environment{
			Core:            coreDiff.removed,
			Remotes:         remotesDiff.removed,
			InstallationType: nextEnv.InstallationType,
		}

		updatedEnv := model.Environment{
			Core:            coreDiff.updated,
			Remotes:         remotesDiff.updated,
			InstallationType: nextEnv.InstallationType,
		}

		if len(coreDiff.added) > 0 || len(remotesDiff.added) > 0 {
			toBeAdded = append(toBeRemoved, addedEnv)
		}

		if len(coreDiff.removed) > 0 || len(remotesDiff.removed) > 0 {
			toBeUpdated = append(toBeRemoved, removedEnv)
		}

		if len(coreDiff.updated) > 0 || len(remotesDiff.updated) > 0 {
			toBeRemoved = append(toBeRemoved, updatedEnv)
		}
	}

	// Check for environments in next that are not in prev (should not happen)
	for key, nextEnv := range nextMap {
		if _, exists := prevMap[key]; !exists {
			toBeRemoved = append(toBeRemoved, nextEnv)
		}
	}

	return
}


func compareCore(prevCore, nextCore map[string]string) diffCore {
	d := diffCore{
		added:   make(map[string]string),
		removed: make(map[string]string),
		updated: make(map[string]string),
	}

	for k, v := range prevCore {
		if nextVal, exists := nextCore[k]; !exists {
			d.removed[k] = v
		} else if v != nextVal {
			d.updated[k] = nextVal
		}
	}

	for k, v := range nextCore {
		if _, exists := prevCore[k]; !exists {
			d.added[k] = v
		}
	}

	return d
}

func compareRemotes(prevRemotes, nextRemotes map[string]map[string]string) diffRemote {
	d := diffRemote{
		added:   make(map[string]map[string]string),
		removed: make(map[string]map[string]string),
		updated: make(map[string]map[string]string),
	}

	for k, prevRemote := range prevRemotes {
		if nextRemote, exists := nextRemotes[k]; !exists {
			d.removed[k] = prevRemote
		} else {
			addedRemote := make(map[string]string)
			updatedRemote := make(map[string]string)
			removedRemote := make(map[string]string)
			for rk, rv := range prevRemote {
				if nextVal, exists := nextRemote[rk]; !exists {
					removedRemote[rk] = rv
				} else if rv != nextVal {
					updatedRemote[rk] = nextVal
				}
			}
			for rk, rv := range nextRemote {
				if _, exists := prevRemote[rk]; !exists {
					addedRemote[rk] = rv
				}
			}
			if len(updatedRemote) > 0 {
				d.updated[k] = updatedRemote
			}
			if len(addedRemote) > 0 {
				d.added[k] = updatedRemote
			}
			if len(updatedRemote) > 0 {
				d.removed[k] = updatedRemote
			}
		}
	}

	for k, nextRemote := range nextRemotes {
		if _, exists := prevRemotes[k]; !exists {
			d.added[k] = nextRemote
		}
	}

	return d
}

// Function to compare Flatpak Applications
func compareApplications(nextRepos []model.Environment, currentApps []model.FlatpakApplication, nextApps []model.FlatpakApplication) ([]model.FlatpakApplication, []model.FlatpakApplication) {
	var appsToAdd []model.FlatpakApplication
	var appsToRemove []model.FlatpakApplication

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

		if !slices.ContainsFunc(nextRepos, func(nextRepo model.Environment) bool {
			return nextRepo.RemoteExists(app.Repo,app.InstallationType)
			//return app.Repo == nextRepo.Name
		}) {
			continue
		}

		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		if _, exists := currentAppMap[key]; !exists {
			appsToAdd = append(appsToAdd, app)
		}
	}

	return appsToAdd, appsToRemove
}

// Function to compare Flatpak Application Permissions (Overrides)
func comparePermissions(currentApps []model.FlatpakApplication, nextApps []model.FlatpakApplication) []model.FlatpakApplication {
	var appsToReplace []model.FlatpakApplication

	for _, nextApp := range nextApps {
		for i, currentApp := range currentApps {
			if nextApp.Name == currentApp.Name && nextApp.Repo == currentApp.Repo && nextApp.InstallationType == currentApp.InstallationType {

				// Compare overrides
				app := model.FlatpakApplication{
					Name:             nextApp.Name,
					Repo:             nextApp.Repo,
					InstallationType: nextApp.InstallationType,
				}

				overridesChanged := false
				for _, value := range nextApp.Overrides {
					if !StringExistsInArray(value, currentApp.Overrides) {
						overridesChanged = true
						app.Overrides = append(app.Overrides, value)
					}
				}
				if overridesChanged {
					currentApps[i].Overrides = nextApp.Overrides
				}

				overridesUserChanged := false
				for _, value := range nextApp.OverridesUser {
					if !StringExistsInArray(value, currentApp.OverridesUser) {
						overridesUserChanged = true
						app.OverridesUser = append(app.OverridesUser, value)
					}
				}
				if overridesUserChanged {
					currentApps[i].OverridesUser = nextApp.OverridesUser
				}

				if overridesChanged || overridesUserChanged {
					appsToReplace = append(appsToReplace, app)
				}
			}
		}
	}

	return appsToReplace
}

func GetDiffState(currentState, nextState model.State) DiffState {
	// Handle differences as needed...
	// Compare repositories
	envToAdd, envToRemove, envToUpdate := compareEnvironments(currentState.Environment, nextState.Environment)
	// Compare applications
	appsToAdd, appsToRemove := compareApplications(nextState.Environment, currentState.Applications, nextState.Applications)
	// Compare permissions
	appsToReplace := comparePermissions(currentState.Applications, nextState.Applications)
	// Create FlatpakDiff structure
	return DiffState{
		EnvToAdd:               envToAdd,
		EnvToRemove:            envToRemove,
		EnvToUpdate:            envToUpdate,
		AppsToAdd:              appsToAdd,
		AppsToRemove:           appsToRemove,
		AppsToReplaceOverrides: appsToReplace,
	}

}
