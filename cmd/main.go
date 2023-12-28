package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"github.com/faan11/flatpak-compose/internal/model"
	"github.com/faan11/flatpak-compose/internal/view"
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
		if (err != nil){
			log.Fatalf("Apply compose file not found: %v \n", err)
			log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
			return
		}

		// Get next state 
		var currentState, nextState model.State
		nextState, err = model.GetFileState(file)
		if (err != nil){
			log.Fatalf("%v \n", err);
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
				view.ExecDiffCommands(diff)
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
		if (err != nil){
			log.Fatalf("Plan compose file not found: %v \n", err)
			log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
			return
		}
		// Get the respective states based on the flag value
		var currentState,nextState model.State
		// Get next state 
		nextState, err = model.GetFileState(file)
		if (err != nil){
			log.Fatalf("%v \n", err);
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
			
			if (err != nil){
				log.Fatalf("Export compose file not found: %v \n", err)
				log.Fatalf("You should specify the input file or create a flatpak-compose.yml or flatpak-compose.yaml in the same directory. \n")
				return
			}
			fmt.Println(file)
			fileState, err := model.GetFileState(file)
			if (err != nil){
				log.Fatalf("%v \n", err);
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
	fmt.Println("Usage:")
	fmt.Println("flatpak-compose apply [-f file.yaml] [-next-state=compose/system-compose]     # Apply changes")
	fmt.Println("flatpak-compose plan [-f file.yaml] [-next-state=compose/system-compose]      # Print commands without applying")
	fmt.Println("flatpak-compose export-state [-f file.yaml] [-state-type=compose/system-compose]  # Export the current state")
	fmt.Println("flatpak-compose help                     # Show help")
}

