package state

import (
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
)

type DiffState struct {
	EnvToAdd               []model.Environment
	EnvToRemove            []model.Environment
	EnvToUpdate            []model.Environment
	AppsToAdd              []model.FlatpakApplication
	AppsToRemove           []model.FlatpakApplication
	PermToAdd 	       []model.FlatpakApplication
	PermToRemove 	       []model.FlatpakApplication
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
			toBeAdded = append(toBeAdded, addedEnv)
		}

		if len(coreDiff.removed) > 0 || len(remotesDiff.removed) > 0 {
			toBeRemoved = append(toBeUpdated, removedEnv)
		}

		if len(coreDiff.updated) > 0 || len(remotesDiff.updated) > 0 {
			toBeUpdated = append(toBeUpdated, updatedEnv)
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
				d.added[k] = addedRemote
			}
			if len(updatedRemote) > 0 {
				d.removed[k] = removedRemote
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

	// currentApps is the system
	// We assume that currentApps is valid.
	// nextApps is the desidered state.

	nextAppMap := make(map[string]bool)
	// Calculate key indexing for the desired state
	for _, app := range nextApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		nextAppMap[key] = true
	}
	// Check if currentApps is not present in the desidered state.
	for _, app := range currentApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		if _, exists := nextAppMap[key]; !exists {
			// if it NOT present, remove it.
			appsToRemove = append(appsToRemove, app)
		}
	}
	// Calculate index of current app.
	currentAppMap := make(map[string]bool)
	for _, app := range currentApps {
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		currentAppMap[key] = true
	}

	// For each application in the desidered state.
	for _, app := range nextApps {
		found:= false;
		// repo Validation?
		for _, nextEnv := range nextRepos {
			if nextEnv.RemoteExists(app.InstallationType,app.Repo){
				found = true;
				break;
			}
		}

		if (!found) {
			fmt.Println("Invalid Remote validation: ",app.Name,  app.Repo, app.InstallationType)
			break;
		}
		// Calculate index key
		key := fmt.Sprintf("%s|%s|%s", app.Repo, app.InstallationType, app.Name)
		// check if the apps exists in the current state.
		if _, exists := currentAppMap[key]; !exists {
			// If it not, add it.
			appsToAdd = append(appsToAdd, app)
		}
	}

	return appsToAdd, appsToRemove
}

// Function to compare Flatpak Application Permissions (Overrides)
func comparePermissions(currentApps []model.FlatpakApplication, nextApps []model.FlatpakApplication) ([]model.FlatpakApplication,[]model.FlatpakApplication)  {
	var appsPermissionRemove,appsPermissionAdd []model.FlatpakApplication
	// Iterate the desidered state (nextApps)
	for _, nextApp := range nextApps {
		// Iterate the current state to find the related couple (nextApp,currentApp)
		for _, currentApp := range currentApps {
			if nextApp.Name == currentApp.Name && nextApp.Repo == currentApp.Repo && nextApp.InstallationType == currentApp.InstallationType {
				// Found it.

				// Compare overrides
				appAdd := model.FlatpakApplication{
					Name:             nextApp.Name,
					Repo:             nextApp.Repo,
					InstallationType: nextApp.InstallationType,
				}

				appRemove := model.FlatpakApplication{
					Name:             nextApp.Name,
					Repo:             nextApp.Repo,
					InstallationType: nextApp.InstallationType,
				}

				//overridesChanged := false
				// Iterate the desidered state.
				for _, value := range nextApp.Overrides {
					// For each key in the desired state, is it available on the previous state?
					if !StringExistsInArray(value, currentApp.Overrides) {
						// NO. need to add it.
						//overridesChanged = true
						appRemove.Overrides = append(appRemove.Overrides, value)
					}
				}
				// Iterate the curent state and see if fields are missing in the desidered state.
				// If yes, please delete it.
				for _, value := range currentApp.Overrides {
					// For each key in the desired state, is it available on the previous state?
					if !StringExistsInArray(value, nextApp.Overrides) {
						// NO. need to add it.
						//overridesChanged = true
						appAdd.Overrides = append(appAdd.Overrides, value)
					}
				}


				for _, value := range nextApp.OverridesUser {
					// For each key in the desired state, is it available on the previous state?
					if !StringExistsInArray(value, currentApp.OverridesUser) {
						// NO. need to add it.
						//overridesChanged = true
						appRemove.OverridesUser = append(appRemove.OverridesUser, value)
					}
				}

				//overridesUserChanged := false
				for _, value := range nextApp.OverridesUser {
					// For each key in the desired state, is it available on the previous state?
					if !StringExistsInArray(value, currentApp.OverridesUser) {
						//overridesUserChanged = true
						appAdd.OverridesUser = append(appAdd.OverridesUser, value)
					}
				}

				if (appAdd.OverridesUser != nil) {
					appsPermissionAdd = append(appsPermissionAdd, appAdd)
				}

				if (appRemove.Overrides != nil) {
					appsPermissionRemove = append(appsPermissionRemove, appRemove)
				}

				// Let's go out... we found the related app.
				break;
			}
		}
	}

	return appsPermissionAdd,appsPermissionRemove
}

func GetDiffState(currentState, nextState model.State) DiffState {
	// Handle differences as needed...
	// Compare repositories
	envToAdd, envToRemove, envToUpdate := compareEnvironments(currentState.Environment, nextState.Environment)
	// Compare applications
	appsToAdd, appsToRemove := compareApplications(nextState.Environment, currentState.Applications, nextState.Applications)
	// Compare permissions
	permToAdd, permToRemove := comparePermissions(currentState.Applications, nextState.Applications)
	// Create FlatpakDiff structure
	return DiffState{
		EnvToAdd:               envToAdd,
		EnvToRemove:            envToRemove,
		EnvToUpdate:            envToUpdate,
		AppsToAdd:              appsToAdd,
		AppsToRemove:           appsToRemove,
		PermToAdd: 		permToAdd,
		PermToRemove: 		permToRemove,
	}

}
