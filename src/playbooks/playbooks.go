package playbooks

import (
	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Playbook []struct {
	Name      string                 `yaml:"name"`
	Groups    interface{}            `yaml:"groups"`
	Variables map[string]interface{} `yaml:"vars"`
	Roles     []struct {
		Role string `yaml:"role"`
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
