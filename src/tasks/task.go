package tasks

import (
	"github.com/smahm006/gear/src/common"
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

// Need custom unmarshal logic for different modules
func (t *Task) UnmarshalYAML(value *yaml.Node) error {
	// Find module name
	t.Tag = value.Tag[1:]
	// De-serialize Task specific fields
	type tempTask Task
	var alias tempTask
	if err := value.Decode(&alias); err != nil {
		return err
	}
	t.Name = alias.Name
	t.With = alias.With
	t.And = alias.And
	// Get module based on tag
	m, err := modules.MapTagToModule(t.Tag, value)
	if err != nil {
		return err
	}
	t.Module = m
	return nil
}

func (t *Task) RunTask(status *common.RunStatus) {
	return
}
