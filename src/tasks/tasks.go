package tasks

import (
	"github.com/smahm006/gear/src/modules"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Tag    string
	Name   string `yaml:"name"`
	Module modules.Module
	With   PreTaskLogic  `yaml:"with"`
	And    PostTaskLogic `yaml:"and"`
}

func NewTasks() *Task {
	return &Task{}
}

// Need custom unmarshal logic for different modules
func (t *Task) UnmarshalYAML(value *yaml.Node) error {
	t.Tag = value.Tag[1:]
	switch t.Tag {
	case "shell":
		var shell_module *modules.ShellModule
		if err := value.Decode(&shell_module); err != nil {
			return err
		}
		t.Module = shell_module
	default:
	}
	type tempTask Task

	var alias tempTask
	if err := value.Decode(&alias); err != nil {
		return err
	}
	t.Name = alias.Name
	t.With = alias.With
	t.And = alias.And
	return nil
}
