package view

import (
	"fmt"
	"github.com/faan11/flatpak-compose/internal/state"
)

func printShellCommands(commands []string) {
	for _, cmd := range commands {
		fmt.Println(cmd)
	}
}

func PrintDiffCommands(diff state.DiffState) {
	list := GenDiffStateCommands(diff)
	printShellCommands(list)
}
