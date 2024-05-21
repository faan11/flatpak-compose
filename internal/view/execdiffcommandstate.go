package view

import (
	"bufio"
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
	"io"
	"os"
	"os/exec"
	"strings"
)

func printOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func executeShellCommandsAndGetOutput(commands []string) {
	for _, cmdStr := range commands {
		fmt.Printf("+ %s \n", cmdStr)
		cmd := exec.Command("sh", "-c", cmdStr)

		// Create pipes to capture stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error creating stdout pipe: %s\n", err)
			continue
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Printf("Error creating stderr pipe: %s\n", err)
			continue
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting command: %s\n", err)
			continue
		}

		// Read stdout and stderr streams concurrently
		go printOutput(stdout)
		go printOutput(stderr)

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Error executing command: %s", err)
		}
	}
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt + " (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		return true
	} else if response == "n" || response == "no" {
		return false
	} else {
		fmt.Println("Please enter y/yes or n/no.")
		return askForConfirmation(prompt)
	}
}

func ExecDiffCommands(diff model.DiffState, assumeyes bool) {
	list := GenDiffStateCommands(diff)
	if len(list) != 0 {
		fmt.Printf("Commands: \n")
		printShellCommands(list)
		if !assumeyes {
			confirmed := askForConfirmation("Are you sure you want to continue?")
			if confirmed {
				fmt.Printf("Execution: \n")
				executeShellCommandsAndGetOutput(list)
				fmt.Println("Completed")
				// Perform the actions you want after confirmation
			} else {
				fmt.Println("Cancelled.")
				// Handle cancellation or exit
			}
		} else {
			fmt.Printf("Execution: \n")
			executeShellCommandsAndGetOutput(list)
			fmt.Println("Completed")
		}
	} else {
		fmt.Println("No changes needs to be done")
	}
}
