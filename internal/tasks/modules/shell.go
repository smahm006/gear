package modules

import (
	"github.com/smahm006/gear/internal/tasks/requonse"
)

type ShellModule struct {
	Cmd         string `yaml:"cmd"`
	Save        string `yaml:"save"`
	FailedWhen  string `yaml:"failed_when"`
	ChangedWhen string `yaml:"changed_when"`
}

func (s *ShellModule) Run() {
}

func (s *ShellModule) Query() requonse.TaskRequest {
	var request requonse.TaskRequest
	return request
}
