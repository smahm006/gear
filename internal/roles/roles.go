package roles

import (
	"fmt"

	"github.com/smahm006/gear/internal/playbook/state"
	"github.com/smahm006/gear/internal/tasks"
	"gopkg.in/yaml.v3"
)

type Role struct {
	Name      string                 `yaml:"name"`
	Variables map[string]interface{} `yaml:"vars"`
	Tasks     []string               `yaml:"tasks"`
	Tags      []string               `yaml:"tags"`
	Handlers  []string               `yaml:"handlers"`
	Path      string
}

func NewRole(name string, variables map[string]interface{}, tags []string) *Role {
	return &Role{
		Name:      name,
		Variables: variables,
		Tags:      tags,
	}
}

func (r *Role) LoadRole() error {
	// Some variables and tags might already exist from the playbook stage
	var temp_role Role
	yaml_data, err := validateRole(r)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yaml_data, &temp_role)
	if err != nil {
		return err
	}
	for k, v := range temp_role.Variables {
		r.Variables[k] = v
	}
	r.Tags = append(r.Tags, temp_role.Tasks...)
	r.Handlers = temp_role.Handlers
	r.Tasks = temp_role.Tasks
	return nil
}

func (r *Role) RunRole(run_state *state.RunState) error {
	var err error
	for _, r_task := range r.Tasks {
		task_path := fmt.Sprintf("%s/tasks/%s", r.Path, r_task)
		tasks := tasks.NewTasks()
		if err = tasks.LoadTasks(task_path); err != nil {
			return err
		}
		if err = tasks.RunTasks(run_state); err != nil {
			return err
		}
	}
	return nil
}
