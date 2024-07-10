package modules

import (
	"github.com/smahm006/gear/internal/connection"
	"github.com/smahm006/gear/internal/tasks/exchange"
)

type ShellModule struct {
	Cmd         string `yaml:"cmd"`
	Save        string `yaml:"save"`
	FailedWhen  string `yaml:"failed_when"`
	ChangedWhen string `yaml:"changed_when"`
}

func (s *ShellModule) Run(request *exchange.TaskRequest, connection connection.Connection) *exchange.TaskResponse {
	var response *exchange.TaskResponse
	response = connection.Execute(s.Cmd)
	return response
}

func (s *ShellModule) Query() *exchange.TaskRequest {
	return &exchange.TaskRequest{Type: exchange.Execute}
}
