package roles

import (
	"fmt"
	"os"

	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/tasks"
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Role struct {
	Name      string                 `yaml:"name"`
	Variables map[string]interface{} `yaml:"vars"`
	Tasks     []string               `yaml:"tasks"`
	Handlers  []string               `yaml:"handlers"`
}

func NewRole() *Role {
	return &Role{}
}

func (r *Role) LoadRole(path string) error {
	yaml_data, err := utils.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yaml_data, &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) RunRole(state *common.RunState) error {
	var err error
	for _, r_task := range r.Tasks {
		wd, _ := os.Getwd()
		task_path := fmt.Sprintf("%s/examples/tasks/%s", wd, r_task)
		fmt.Println(task_path)
		tasks := tasks.NewTasks()
		if err = tasks.LoadTasks(task_path); err != nil {
			return err
		}
		if err = tasks.RunTasks(state); err != nil {
			return err
		}
	}
	return nil
}
