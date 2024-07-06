package playbook

import (
	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/inventory"
	"github.com/smahm006/gear/src/roles"
	"github.com/smahm006/gear/src/tasks"
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Playbook []Play

type Play struct {
	Name      string                 `yaml:"name"`
	Groups    interface{}            `yaml:"groups"`
	Variables map[string]interface{} `yaml:"vars"`
	PreTasks  []tasks.Task           `yaml:"pre"`
	PostTasks []tasks.Task           `yaml:"post"`
	Roles     []struct {
		Name      string                 `yaml:"role"`
		Variables map[string]interface{} `yaml:"vars"`
		Tags      []string               `yaml:"tags"`
	} `yaml:"roles"`
}

func NewPlaybook() *Playbook {
	return &Playbook{}
}

func (p *Playbook) LoadPlaybook(path string) error {
	yaml_data, err := utils.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yaml_data, &p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Playbook) RunPlaybook(cli *cmd.CliParser, i *inventory.Inventory) error {
	var err error
	state := common.NewRunState(cli, i)
	for _, play := range *p {
		for _, pre_task := range play.PreTasks {
			_ = pre_task
		}
		for _, p_role := range play.Roles {
			role := roles.NewRole(p_role.Name, p_role.Variables, p_role.Tags)
			if err = role.LoadRole(); err != nil {
				return err
			}
			if err = role.RunRole(state); err != nil {
				return err
			}
		}
		for _, post_task := range play.PostTasks {
			_ = post_task
		}
	}
	return nil
}
