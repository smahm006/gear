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

type Play struct {
	Name      string                 `yaml:"name"`
	Groups    []string               `yaml:"groups"`
	Variables map[string]interface{} `yaml:"vars"`
	PreTasks  []tasks.Task           `yaml:"pre"`
	PostTasks []tasks.Task           `yaml:"post"`
	Roles     []struct {
		Name      string                 `yaml:"role"`
		Variables map[string]interface{} `yaml:"vars"`
		Tags      []string               `yaml:"tags"`
	} `yaml:"roles"`
}

type Playbook []Play

func NewPlaybook() *Playbook {
	return &Playbook{}
}

func (p *Playbook) LoadPlaybook(cli *cmd.CliParser, inventory *inventory.Inventory) error {
	yaml_data, err := utils.ReadFile(cli.PlaybookPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yaml_data, &p)
	if err != nil {
		return err
	}
	for i := range *p {
		var err error
		play := &(*p)[i]
		if err = validateGroups(inventory, play); err != nil {
			return err
		}
	}
	return nil
}

func (p *Playbook) RunPlaybook(cli *cmd.CliParser, i *inventory.Inventory) error {
	var err error
	state := common.NewRunState(cli, i)
	for i := range *p {
		play := &(*p)[i]
		collected_hosts := collectHosts(state, play)
		collected_vars := collectVars(state, play)
		status := common.NewRunStatus(collected_hosts, collected_vars)
		for _, pre_task := range play.PreTasks {
			pre_task.RunTask(status)
		}
		for _, p_role := range play.Roles {
			role := roles.NewRole(p_role.Name, p_role.Variables, p_role.Tags)
			if err = role.LoadRole(); err != nil {
				return err
			}
			if err = role.RunRole(status); err != nil {
				return err
			}
		}
		for _, post_task := range play.PostTasks {
			post_task.RunTask(status)
		}
	}
	return nil
}
