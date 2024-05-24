package main

import (
	"flag"
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
	"github.com/faan11/flatpak-compose/internal/view"
	"log"
	"os"
)

func getValidFileName(defaultFileName string) (string, error) {
	fileExists := func(name string) bool {
		_, err := os.Stat(name)
		return err == nil
	}
	fileNames := []string{defaultFileName, "flatpak-compose.yml"} // Add more file names if needed
	for _, fileName := range fileNames {
		if fileName != "" {
			if fileExists(fileName) {
				return fileName, nil
			} else {
				return "", fmt.Errorf("No valid input file found")
			}
		}
	}

	return "", fmt.Errorf("No valid input file found")
}

func main() {
	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	applyFile := applyCmd.String("f", "flatpak-compose.yaml", "YAML file for applying changes")
	applyNextState := applyCmd.String("current-state", "system-compose", "Specify the current state type: system-compose or system")
	applyAssumeyes := applyCmd.Bool("assumeyes", false, "Automatically answer yes for all questions")

	planCmd := flag.NewFlagSet("plan", flag.ExitOnError)
	planFile := planCmd.String("f", "flatpak-compose.yaml", "YAML file for planning changes")
	planNextState := planCmd.String("current-state", "system-compose", "Specify the current state type: system-compose or system")

	exportCmd := flag.NewFlagSet("export-state", flag.ExitOnError)
	exportFile := exportCmd.String("f", "flatpak-compose.yaml", "YAML file for exporting state")

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "apply":
		applyCmd.Parse(os.Args[2:])
		// Check if next-state is valid
		if *applyNextState != "system-compose" && *applyNextState != "system" {
			log.Fatalf("Invalid next-state type. Use 'system-compose' or 'system'.")
			return
		}

		file, err := getValidFileName(*applyFile)
		if err != nil {
			log.Fatalf("Apply compose file not found: %v \n", err)
			log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
			return
		}

		// Get next state
		var currentState, nextState model.State
		nextState, err = model.GetFileState(file)
		if err != nil {
			log.Fatalf("%v \n", err)
			return
		}

		switch *applyNextState {
		case "system-compose":
			currentState = model.GetComposeState(nextState, model.GetSystemState())
		case "system":
			currentState = model.GetSystemState()
		}

		diff := model.GetDiffState(currentState, nextState)

		if applyCmd.Parsed() {
			if planCmd.Parsed() {
				view.PrintDiffCommands(diff)
			} else {
				view.ExecDiffCommands(diff, *applyAssumeyes)
			}
		}

	case "plan":
		planCmd.Parse(os.Args[2:])
		// Check if next-state is valid
		if *planNextState != "system-compose" && *planNextState != "system" {
			log.Fatalf("Invalid next-state type. Use 'compose' or 'system'.")
			return
		}
		// Get valid file
		file, err := getValidFileName(*planFile)
		if err != nil {
			log.Fatalf("Plan compose file not found: %v \n", err)
			log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
			return
		}
		// Get the respective states based on the flag value
		var currentState, nextState model.State
		// Get next state
		nextState, err = model.GetFileState(file)
		if err != nil {
			log.Fatalf("%v \n", err)
			return
		}

		switch *applyNextState {
		case "system-compose":
			currentState = model.GetComposeState(nextState, model.GetSystemState())
		case "system":
			currentState = model.GetSystemState()
		}

		diff := model.GetDiffState(currentState, nextState)

		view.PrintDiffCommands(diff)

	case "export-state":
		// Check if state-type is valid

		if len(os.Args) < 3 {
			log.Fatal("Specify the state type to export: 'current-compose' or 'system'")
		}
		exportCmd.Parse(os.Args[3:])

		exportStateType := os.Args[2]

		if exportStateType != "system-compose" && exportStateType != "system" {
			log.Fatalf("Invalid state-type to export. Use 'compose' or 'system'.")
			return
		}

		// Get the respective state based on the flag value
		var exportState model.State
		switch exportStateType {
		case "system-compose":
			// Get next state file name.
			file, err := getValidFileName(*exportFile)

			if err != nil {
				log.Fatalf("Export compose file not found: %v \n", err)
				log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
				return
			}
			fmt.Println(file)
			fileState, err := model.GetFileState(file)
			if err != nil {
				log.Fatalf("%v \n", err)
				return
			}
			exportState = model.GetComposeState(fileState, model.GetSystemState())
		case "system":
			exportState = model.GetSystemState()
		}
		// Export the state to the file
		// Replace this with the actual export logic
		view.PrintState(exportState)
	case "help":
		printUsage()

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Flatpak-Compose: A utility tool for managing Flatpak applications and repositories")
	fmt.Println("\nFlatpak-Compose is a command-line utility designed to facilitate managing Flatpak applications and repositories. It provides functionalities to apply, plan, and export states of applications and repositories in Flatpak. With this tool, users can identify changes between the system setup and the compose file, enabling efficient management and deployment of Flatpak applications.")
	fmt.Println("The tool performs the difference between the current and the compose state. The current state can be the system state or the intersection between the system and compose state (system-compose).")
	fmt.Println("The current state is the system-compose state by default in order to avoid unwanted changes.")
	fmt.Println("\nUsage:")
	fmt.Println("flatpak-compose apply [-f file.yaml] [-current-state=system/system-compose]     # Apply changes based on the difference between the current state and the system state")
	fmt.Println("flatpak-compose plan [-f file.yaml] [-current-state=system/system-compose]      # Show changes based on the difference between the current state and the system state")
	fmt.Println("flatpak-compose export-state [-f file.yaml] [-state-type=system/system-compose]  # Show the system or system-compose state using the YAML format")
	fmt.Println("\nOptions:")
	fmt.Println("  apply         : Apply changes based on the difference between the current state and the system state")
	fmt.Println("  plan          : Show changes based on the difference between the current state and the system state")
	fmt.Println("  export-state  : Show the system or system-compose state using the YAML format")
	fmt.Println("  help          : Show usage information")
	fmt.Println("\nFlags:")
	fmt.Println("  -f                : YAML file to load (default: flatpak-compose.yaml)")
	fmt.Println("  -current-state    : Specify the current state type (system/system-compose)")
	fmt.Println("  -state-type       : Specify the state type to export (system/system-compose)")
	fmt.Println("\nExplanation:")
	fmt.Println("  current state     : Can be the system or system-compose state")
	fmt.Println("  system state      : Includes all the applications/repos in the system")
	fmt.Println("  system-compose state: Includes all the application/repos that are in common between the compose and the system (right join)")
}
