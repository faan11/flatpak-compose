package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/faan11/flatpak-compose/internal/model"
	"github.com/faan11/flatpak-compose/internal/view"
)

func main() {
	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	applyFile := applyCmd.String("f", "flatpak-compose.yaml", "YAML file for applying changes")

	planCmd := flag.NewFlagSet("plan", flag.ExitOnError)
	planFile := planCmd.String("f", "flatpak-compose.yaml", "YAML file for planning changes")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "apply":
		applyCmd.Parse(os.Args[2:])

		nextState, err := model.GetNextState(*applyFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		currentState := model.GetCurrentState()

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

		nextState, err := model.GetNextState(*planFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		currentState := model.GetCurrentState()

		diff := model.GetDiffState(currentState, nextState)
		view.PrintDiffCommands(diff)
	
	case "export":
		exportCmd.Parse(os.Args[2:])

		currentState := model.GetCurrentState()
		view.PrintState(currentState)
	case "help":
		printUsage()

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("flatpak-compose apply [-f file.yaml]     # Apply changes")
	fmt.Println("flatpak-compose plan [-f file.yaml]      # Print commands without applying")
	fmt.Println("flatpak-compose export      	      # Print the current state in a yaml file")
	fmt.Println("flatpak-compose help                     # Show help")
}

