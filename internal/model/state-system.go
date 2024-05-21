package model

import (
	"github.com/faan11/flatpak-compose/internal/utility"
	"log"
	"os/exec"
	"strings"
)

func GetSystemState() State {
	var currentState State

	// Get list of installed applications
	installedAppsCmd := exec.Command("flatpak", "list", "--app", "--columns=application,branch,origin,installation")
	installedAppsOutput, err := installedAppsCmd.Output()
	if err != nil {
		log.Fatalf("Error getting installed applications: %s\n", err)
	}

	// Parse installed applications output
	installedApps := strings.Split(string(installedAppsOutput), "\n")
	for _, app := range installedApps {
		fields := strings.Fields(app)
		if len(fields) >= 2 {
			currentState.Applications = append(currentState.Applications, FlatpakApplication{
				Name:             fields[0],
				Branch:           fields[1],
				Repo:             fields[2],
				InstallationType: fields[3],
				// Add other properties as needed
			})
		}
	}

	// Get list of remotes
	remotesCmd := exec.Command("flatpak", "remote-list", "--columns=name,url,options")
	remotesOutput, err := remotesCmd.Output()
	if err != nil {
		log.Fatalf("Error getting remotes: %s\n", err)
	}

	// Parse remotes output
	remotes := strings.Split(string(remotesOutput), "\n")
	for _, remote := range remotes {
		fields := strings.Fields(remote)
		if len(fields) >= 2 {
			currentState.Repos = append(currentState.Repos, FlatpakRepo{
				Name:             fields[0],
				URI:              fields[1],
				InstallationType: fields[2],
			})
		}
	}

	// Get permissions (overrides) for installed applications
	for i, app := range currentState.Applications {
		permissionsCmd := exec.Command("flatpak", "override", app.Name, "--show")
		permissionsOutput, err := permissionsCmd.Output()
		if err != nil {
			log.Fatalf("Error getting permissions for %s: %s\n", app.Name, err)
		}

		currentState.Applications[i].Overrides = utility.MapPermissionsToFlatpakOverrideFlags(string(permissionsOutput))
	}

	// Get permissions (overrides) for installed applications
	for i, app := range currentState.Applications {
		permissionsCmd := exec.Command("flatpak", "override", app.Name, "--show", "--user")
		permissionsOutput, err := permissionsCmd.Output()
		if err != nil {
			log.Fatalf("Error getting permissions for %s: %s\n", app.Name, err)
		}

		currentState.Applications[i].OverridesUser = utility.MapPermissionsToFlatpakOverrideFlags(string(permissionsOutput))
	}

	// Get permissions (all) for installed applications
	for i, app := range currentState.Applications {
		permissionsCmd := exec.Command("flatpak", "info", app.Name, "-M")
		permissionsOutput, err := permissionsCmd.Output()
		if err != nil {
			log.Fatalf("Error getting permissions for %s: %s\n", app.Name, err)
		}

		currentState.Applications[i].All = utility.MapPermissionsToFlatpakOverrideFlags(string(permissionsOutput))
	}
	return currentState
}
