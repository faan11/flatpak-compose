package view

import (
	"fmt"
	"github.com/faan11/flatpak-compose/internal/model"
	"gopkg.in/yaml.v2"
)

func PrintState(state model.State) {
	stateYAML, err := yaml.Marshal(state)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(stateYAML))
}
