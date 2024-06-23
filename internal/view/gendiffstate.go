package view

import (
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"github.com/faan11/flatpak-compose/internal/model"
	"github.com/faan11/flatpak-compose/internal/state"
	"github.com/faan11/flatpak-compose/internal/utility"
)


// ConvertMapToText converts a map to a multiline text string
func ConvertMapToOptions(m map[string]string) string {
	var result string

	if title, ok := m["xa.title"]; ok {
		result += fmt.Sprintf("--title=%s\n", title)
	}

	if url, ok := m["url"]; ok {
		result += fmt.Sprintf("--url=%s\n", url)
	}

	if homepage, ok := m["xa.homepage"]; ok {
		result += fmt.Sprintf("--homepage=%s\n", homepage)
	}

	if comment, ok := m["xa.comment"]; ok {
		result += fmt.Sprintf("--comment=%s\n", comment)
	}

	if description, ok := m["xa.description"]; ok {
		result += fmt.Sprintf("--description=%s\n", description)
	}

	if icon, ok := m["xa.icon"]; ok {
		result += fmt.Sprintf("--icon=%s\n", icon)
	}

	if gpgKey, ok := m["GPGKey"]; ok {
		file, errs := ioutil.TempFile("", "*.gpg")
		if errs != nil {
		   fmt.Println(errs)
		   return result 
		}
		// Decode the base64 string
		decodedBytes, err := base64.StdEncoding.DecodeString(gpgKey)
		if err != nil {
			fmt.Println("Error decoding base64 string:", err)
			return result
		}
		// Write the binary data to the file
		_, err = file.Write(decodedBytes)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return result
		}
		// Close the file
		errs = file.Close()
		if errs != nil {
		   fmt.Println(errs)
		   return result
		}
		result += fmt.Sprintf("--gpg-import=%s\n", file.Name())
	}

	return result
}

// ConvertMapToText converts a map to a multiline text string
func ConvertMapToText(m map[string]string) string {
	var result string

	result += "[Flatpak Repo]\n"

	if title, ok := m["xa.title"]; ok {
		result += fmt.Sprintf("Title=%s\n", title)
	}

	if url, ok := m["url"]; ok {
		result += fmt.Sprintf("Url=%s\n", url)
	}

	if homepage, ok := m["xa.homepage"]; ok {
		result += fmt.Sprintf("Homepage=%s\n", homepage)
	}

	if comment, ok := m["xa.comment"]; ok {
		result += fmt.Sprintf("Comment=%s\n", comment)
	}

	if description, ok := m["xa.description"]; ok {
		result += fmt.Sprintf("Description=%s\n", description)
	}

	if icon, ok := m["xa.icon"]; ok {
		result += fmt.Sprintf("Icon=%s\n", icon)
	}

	if gpgKey, ok := m["GPGKey"]; ok {
		result += fmt.Sprintf("GPGKey=%s\n", gpgKey)
	}

	return result
}

// Function to generate Flatpak commands to add repositories
func generateEnvAddCommands(envs []model.Environment) []string {
	var commands []string

	for _, env := range envs {
		for name, remote := range env.Remotes {
			// Here we need to create a temporary remote flatpakrepo file.
			// Create a temporary file but don't delete it
			file, errs := ioutil.TempFile("", "*.flatpakrepo")
			if errs != nil {
			   fmt.Println(errs)
			   return []string{} 
			}
			text := ConvertMapToText(remote)
			// Write some text to the file
			_, errs = file.WriteString(text)
			if errs != nil {
			   fmt.Println(errs)
			   return []string{} 
			}
	        	// Close the file
			errs = file.Close()
			if errs != nil {
			   fmt.Println(errs)
			   return []string{} 
			}
		

			cmd := fmt.Sprintf("flatpak remote-add --%s --if-not-exists %s file://%s ",  env.InstallationType, name, file.Name())
			// Adds no verification if it is needed.
			if verify, ok := remote["gpg-verify"]; ok && verify == "false" {
				cmd += "--no-gpg-verify"
			}

			commands = append(commands, cmd)
		}
	}

	return commands
}

// Function to generate Flatpak commands to remove repositories
func generateEnvRemoveCommands(envs []model.Environment) []string {
	var commands []string

	for _, env := range envs {
		for name, _ := range env.Remotes {
			cmd := fmt.Sprintf("flatpak remote-delete --%s %s", env.InstallationType, name)
			commands = append(commands, cmd)
		}
	}

	return commands
}

// Function to generate Flatpak commands to update repositories
func generateEnvUpdateCommands(envs []model.Environment) []string {
	var commands []string

	for _, env := range envs {
		for name, remote := range env.Remotes {
			options := ConvertMapToOptions(remote)
			cmd := fmt.Sprintf("flatpak remote-modify %s --%s %s ", options, env.InstallationType, name)
			commands = append(commands, cmd)
		}
	}

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
func generateAppPermissionsCommands(added, removed []model.FlatpakApplication) []string {
	var commands []string

	for _, app := range added {
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
	
	for _, app := range removed {
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
				cmd += fmt.Sprintf("%s ", utility.NegateFlag(value))
			}
			commands = append(commands, cmd)
		}
	}

	return commands
}

// Function to generate Flatpak commands to replace dynamic permissions
func generateAppDynamicPermissionsCommands(added, removed []model.FlatpakApplication) []string {
	var commands []string
	// Order is important. to apply changes....
	// Delete =remove
	// Add = add
	// Change = remove, add
	// First remove	
	for _, app := range removed {
		for _, p := range app.Permissions {
			cmd := fmt.Sprintf("flatpak permission-remove \"%s\" \"%s\" \"%s\" ", p.Table, p.Object, app.Name)
			commands = append(commands, cmd)
		}
	}
	// then add
	for _, app := range added {
		for _, p := range app.Permissions {
			cmd := fmt.Sprintf("flatpak permission-set --data=\"%s\" \"%s\" \"%s\" \"%s\" \"%s\" ", p.Data, p.Table, p.Object, app.Name, p.Permission)
			commands = append(commands, cmd)
		}
	}
	return commands
}

func GenDiffStateCommands(diff state.DiffState) []string {
	var commands []string

	envRemoveCommands := generateEnvRemoveCommands(diff.EnvToRemove)
	commands = append(commands, envRemoveCommands...)

	// Generate commands for repositories
	envAddCommands := generateEnvAddCommands(diff.EnvToAdd)
	commands = append(commands, envAddCommands...)

	// Generate commands for repositories
	envUpdateCommands := generateEnvUpdateCommands(diff.EnvToUpdate)
	commands = append(commands, envUpdateCommands...)

	appUninstallCommands := generateAppUninstallCommands(diff.AppsToRemove)
	commands = append(commands, appUninstallCommands...)

	// Generate commands for applications
	appInstallCommands := generateAppInstallCommands(diff.AppsToAdd)
	commands = append(commands, appInstallCommands...)

	// Generate commands for replacing permissions
	permAddRemoveCommands := generateAppPermissionsCommands(diff.PermToAdd,diff.PermToRemove)
	commands = append(commands, permAddRemoveCommands...)

	// Generate commands for replacing permissions
	dynamicPermAddRemoveCommands := generateAppDynamicPermissionsCommands(diff.DynamicPermToAdd,diff.DynamicPermToRemove)
	commands = append(commands, dynamicPermAddRemoveCommands...)
	
	return commands
}
