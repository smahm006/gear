package playbook

import (
	"fmt"
	"os"

	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/inventory"
	"github.com/smahm006/gear/src/roles"
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Playbook []Play

type Play struct {
	Name      string                 `yaml:"name"`
	Groups    interface{}            `yaml:"groups"`
	Variables map[string]interface{} `yaml:"vars"`
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
		for _, p_role := range play.Roles {
			wd, _ := os.Getwd()
			role_path := fmt.Sprintf("%s/examples/roles/%s", wd, p_role.Name)
			role := roles.NewRole()
			if err = role.LoadRole(role_path); err != nil {
				return err
			}
			if err = role.RunRole(state); err != nil {
				return err
			}
		}
	}
	return nil
}
