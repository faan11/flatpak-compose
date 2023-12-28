package view;
import (
	"fmt"
	"os/exec"
	"io"
	"bufio"
	"github.com/faan11/flatpak-compose/internal/model"
)



func printOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func executeShellCommandsAndGetOutput(commands []string) {
	for _, cmdStr := range commands {
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
			fmt.Printf("Error executing command: %s\n", err)
		}
	}
}


func ExecDiffCommands(diff model.DiffState) {
	list := GenDiffStateCommands(diff)
	executeShellCommandsAndGetOutput(list)
}
