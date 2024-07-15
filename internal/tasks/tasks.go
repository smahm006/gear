package tasks

import (
	"fmt"

	"github.com/smahm006/gear/internal/playbook/state"
	"github.com/smahm006/gear/internal/tasks/modules"
	"github.com/smahm006/gear/internal/utils"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Tag    string
	Name   string `yaml:"name"`
	Module modules.Module
	With   *modules.ModuleWith `yaml:"with"`
	And    *modules.ModuleAnd  `yaml:"and"`
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
	// Set default items so we always iterate over one
	if t.With == nil {
		t.With = &modules.ModuleWith{
			Items: []interface{}{""},
		}
	} else if t.With.Items == nil {
		t.With.Items = []interface{}{""}
	}
	t.And = alias.And
	// Get module based on tag
	m, err := modules.MapTagToModule(t.Tag, value)
	if err != nil {
		return err
	}
	t.Module = m
	return nil
}

type Tasks []*Task

func NewTasks() *Tasks {
	return &Tasks{}
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

func (t *Tasks) RunTasks(run_state *state.RunState) error {
	collected_tasks, err := collectTasks(run_state, *t)
	if err != nil {
		return err
	}
	for _, task := range collected_tasks {
		for _, item := range task.With.Items {
			if err := task.RunTask(run_state, item); err != nil {
				fmt.Println("ERROR: ", err)
			}
		}
	}
	return nil
}
