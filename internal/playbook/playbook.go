package playbook

import (
	"github.com/smahm006/gear/internal/cli"
	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/playbook/state"
	"github.com/smahm006/gear/internal/roles"
	"github.com/smahm006/gear/internal/tasks"
	"github.com/smahm006/gear/internal/utils"
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

func (p *Playbook) LoadPlaybook(cli *cli.CliParser, inventory *inventory.Inventory) error {
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

func (p *Playbook) RunPlaybook(cli *cli.CliParser, i *inventory.Inventory) error {
	run_state := state.NewRunState(cli, i)
	for index := range *p {
		play := &(*p)[index]
		collected_hosts, err := collectHosts(run_state, play)
		if err != nil {
			return err
		}
		if err = validateHosts(collected_hosts, play); err != nil {
			return err
		}
		run_status := state.NewRunStatus(collected_hosts, play.Variables)
		run_state.Status = run_status
		for _, pre_task := range play.PreTasks {
			for _, item := range pre_task.With.Items {
				pre_task.RunTask(run_state, item)
			}
		}
		for _, p_role := range play.Roles {
			role := roles.NewRole(p_role.Name, p_role.Variables, p_role.Tags)
			if err = role.LoadRole(); err != nil {
				return err
			}
			if err = role.RunRole(run_state); err != nil {
				return err
			}
		}
		for _, post_task := range play.PostTasks {
			for _, item := range post_task.With.Items {
				post_task.RunTask(run_state, item)
			}
		}
		CleanUpPlay(run_state.Status)
	}
	return nil
}
