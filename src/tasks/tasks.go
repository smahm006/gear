package tasks

import (
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Tasks []Task

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

func (t *Tasks) RunTasks(state *common.RunState) error {
	for _, task := range *t {
		utils.PrintMap(task)
		task.Module.Run()
	}
	return nil
}
