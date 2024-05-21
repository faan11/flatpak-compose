package view

import (
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
)

// Function to generate Flatpak commands to add repositories
func generateRepoAddCommands(repos []model.FlatpakRepo) []string {
	var commands []string

	/*for _, repo := range repos {
		cmd := fmt.Sprintf("flatpak remote-add --%s --if-not-exists %s %s",  repo.InstallationType, repo.Name, repo.URI)
		commands = append(commands, cmd)
	}*/

	return commands
}

// Function to generate Flatpak commands to remove repositories
func generateRepoRemoveCommands(repos []model.FlatpakRepo) []string {
	var commands []string

	/*for _, repo := range repos {
		cmd := fmt.Sprintf("flatpak remote-delete --%s %s", repo.InstallationType, repo.Name)
		commands = append(commands, cmd)
	}*/

	return commands
}

// Function to generate Flatpak commands to install applications
func generateAppInstallCommands(apps []model.FlatpakApplication) []string {
	var commands []string

	for _, app := range apps {
		cmd := fmt.Sprintf("flatpak install %s %s --%s --assumeyes", app.Repo, app.Name, app.InstallationType)
		commands = append(commands, cmd)
		// Adds permissions if exists
		if len(app.Overrides) != 0 {
			cmd = fmt.Sprintf("flatpak override --system %s ", app.Name)
			for _, value := range app.Overrides {
				cmd += fmt.Sprintf("%s ", value)
			}
			commands = append(commands, cmd)
		}
		if len(app.OverridesUser) != 0 {
			cmd = fmt.Sprintf("flatpak override --user %s ", app.Name)
			for _, value := range app.OverridesUser {
				cmd += fmt.Sprintf("%s ", value)
			}
			commands = append(commands, cmd)
		}
	}

	return commands
}

// Function to generate Flatpak commands to uninstall applications
func generateAppUninstallCommands(apps []model.FlatpakApplication) []string {
	var commands []string

	for _, app := range apps {
		cmd := fmt.Sprintf("flatpak uninstall --%s --assumeyes %s", app.InstallationType, app.Name)
		commands = append(commands, cmd)
	}

	return commands
}

// Function to generate Flatpak commands to replace permissions (overrides)
func generateAppReplacePermissionsCommands(apps []model.FlatpakApplication) []string {
	var commands []string

	for _, app := range apps {
		if len(app.Overrides) != 0 {
			cmd := fmt.Sprintf("flatpak override --system %s ", app.Name)
			for _, value := range app.Overrides {
				cmd += fmt.Sprintf("%s ", value)
			}
			commands = append(commands, cmd)
		}
		if len(app.OverridesUser) != 0 {
			cmd := fmt.Sprintf("flatpak override --user %s ", app.Name)
			for _, value := range app.OverridesUser {
				cmd += fmt.Sprintf("%s ", value)
			}
			commands = append(commands, cmd)
		}
	}

	return commands
}

func GenDiffStateCommands(diff model.DiffState) []string {
	var commands []string

	repoRemoveCommands := generateRepoRemoveCommands(diff.ReposToRemove)
	commands = append(commands, repoRemoveCommands...)

	// Generate commands for repositories
	repoAddCommands := generateRepoAddCommands(diff.ReposToAdd)
	commands = append(commands, repoAddCommands...)

	appUninstallCommands := generateAppUninstallCommands(diff.AppsToRemove)
	commands = append(commands, appUninstallCommands...)

	// Generate commands for applications
	appInstallCommands := generateAppInstallCommands(diff.AppsToAdd)
	commands = append(commands, appInstallCommands...)

	// Generate commands for replacing permissions
	appReplacePermissionsCommands := generateAppReplacePermissionsCommands(diff.AppsToReplaceOverrides)
	commands = append(commands, appReplacePermissionsCommands...)

	return commands
}
