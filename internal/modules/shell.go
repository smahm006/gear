package modules

import "github.com/smahm006/gear/internal/state"

type ShellModule struct {
	Cmd         string `yaml:"cmd"`
	Save        string `yaml:"save"`
	FailedWhen  string `yaml:"failed_when"`
	ChangedWhen string `yaml:"changed_when"`
}

func (s *ShellModule) Run() {
}

func (s *ShellModule) Query() state.TaskRequest {
	var request state.TaskRequest
	return request
}
