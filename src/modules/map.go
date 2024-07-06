package modules

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func MapTagToModule(tag string, value *yaml.Node) (Module, error) {
	var m Module
	switch tag {
	case "shell":
		var shell_module *ShellModule
		if err := value.Decode(&shell_module); err != nil {
			return m, err
		}
		m = shell_module
	default:
		return m, fmt.Errorf("Module %s not found", tag)
	}
	return m, nil
}
