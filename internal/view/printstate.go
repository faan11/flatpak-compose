package view;

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/faan11/flatpak-compose/internal/model"
)

func PrintState(state model.State) {
	stateYAML, err := yaml.Marshal(state)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(stateYAML))
}

