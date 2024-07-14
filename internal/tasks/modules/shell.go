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

func (s *ShellModule) Run(connection connection.Connection, request *exchange.TaskRequest, with *ModuleWith, and *ModuleAnd) *exchange.TaskResponse {
	var response *exchange.TaskResponse
	switch request.Type {
	case exchange.Execute:
		response = s.Execute(connection, request, with, and)
		response.Type = exchange.Executed
	default:
		return response
	}

	return response
}

func (s *ShellModule) Query() *exchange.TaskRequest {
	return &exchange.TaskRequest{Type: exchange.Execute}
}

func (s *ShellModule) Execute(connection connection.Connection, request *exchange.TaskRequest, with *ModuleWith, and *ModuleAnd) *exchange.TaskResponse {
	return connection.Execute(s.Cmd)
}
