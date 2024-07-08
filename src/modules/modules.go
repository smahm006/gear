package modules

import (
	"fmt"

	"github.com/smahm006/gear/src/common"
	"gopkg.in/yaml.v3"
)

type Module interface {
	Query() common.TaskRequest
	Run()
}

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
