package view

import (
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
)

func printShellCommands(commands []string) {
	for _, cmd := range commands {
		fmt.Println(cmd)
	}
}

func PrintDiffCommands(diff model.DiffState) {
	list := GenDiffStateCommands(diff)
	printShellCommands(list)
}
