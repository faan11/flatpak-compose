package state 

import (
	"log"
	"os/exec"
	"bufio"
	"fmt"
	"strings"
	"github.com/faan11/flatpak-compose/internal/utility"
	"github.com/faan11/flatpak-compose/internal/model"
)

func GetSystemState() model.State {
	var currentState model.State

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
			currentState.Applications = append(currentState.Applications, model.FlatpakApplication{
				Name:             fields[0],
				Branch:           fields[1],
				Repo:             fields[2],
				InstallationType: fields[3],
				// Add other properties as needed
			})
		}
	}

	//
	// Get user environment 
	//
	uEnv, err := utility.GetUserEnvironment()
	if err != nil {
		log.Fatalf("Error getting user env: %v\n", err)
	}
	currentState.Environment = append(currentState.Environment, uEnv)
	//
	// Get system environment
	//
	sEnv, err := utility.GetSystemEnvironment()
	if err != nil {
		log.Fatalf("Error getting system env: %v\n", err)
	}
	currentState.Environment = append(currentState.Environment, sEnv)


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


	for i, app := range currentState.Applications {
		permissionsCmd := exec.Command("flatpak", "permission-show", app.Name)
		stdout, err := permissionsCmd.StdoutPipe()
		if err != nil {
		    log.Fatalf("Error creating stdout pipe:", err)
		}
		
		if err := permissionsCmd.Start(); err != nil {
		    log.Fatalf("Error starting command:", err)
		}

		scanner := bufio.NewScanner(stdout)

	    	// Skip the first line
	    	if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
			    log.Fatalf("Error scanning permission output:", err)
			}
		}

		for scanner.Scan() {
		    line := scanner.Text()
		    permission, err := parsePermission(line)
		    if err != nil {
			log.Fatalf("Error parsing permission:", err)
			continue
		    }
		    currentState.Applications[i].Permissions = append(currentState.Applications[i].Permissions, permission)
		}

		if err := scanner.Err(); err != nil {
		    log.Fatalf("Error scanning permission output:", err)
		}

		if err := permissionsCmd.Wait(); err != nil {
		    log.Fatalf("Error waiting for command:", err)
		}
    	}

	return currentState
}


func parsePermission(line string) (model.Permission, error) {
    /*parts := strings.Fields(line)
    if len(parts) != 5 {
        return model.Permission{}, fmt.Errorf("line does not have the correct number of fields: %d \"%s\"", len(parts), line)
    }*/
    parts := strings.SplitN(line, "\t", 5) // Split on tab, maximum 5 splits

    if len(parts) < 5 {
        return model.Permission{}, fmt.Errorf("line does not have the correct number of fields: %d", len(parts))
    }
    // Combine the additional fields into the last field
    return model.Permission{
        Table:       parts[0],
        Object:      parts[1],
        //App:         parts[2],
        Permission: parts[3],
        Data:        parts[4],
    }, nil
}
