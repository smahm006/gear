package tasks

import (
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/modules"
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Tasks []Task

type Task struct {
	Tag    string
	Name   string `yaml:"name"`
	Module modules.Module
	With   PreTaskLogic  `yaml:"with"`
	And    PostTaskLogic `yaml:"and"`
}

func NewTasks() *Tasks {
	return &Tasks{}
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

func (t *Tasks) LoadTasks(path string) error {
	yaml_data, err := utils.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yaml_data, &t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tasks) RunTasks(state *common.RunState) error {
	for _, task := range *t {
		utils.PrintMap(task)
		task.Module.Run()
	}
	return nil
}
