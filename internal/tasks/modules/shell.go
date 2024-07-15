package modules

import (
	"github.com/smahm006/gear/internal/connection"
	"github.com/smahm006/gear/internal/tasks/exchange"
	"github.com/smahm006/gear/internal/templar"
)

type ShellModule struct {
	Cmd         string `yaml:"cmd"`
	Save        string `yaml:"save"`
	FailedWhen  string `yaml:"failed_when"`
	ChangedWhen string `yaml:"changed_when"`
}

func (s *ShellModule) Run(connection connection.Connection, request *exchange.TaskRequest, vars map[string]interface{}) *exchange.TaskResponse {
	var response *exchange.TaskResponse
	switch request.Type {
	case exchange.Execute:
		response = s.Execute(connection, request, vars)
		response.Type = exchange.Executed
	default:
		return response
	}

	return response
}

func (s *ShellModule) Query() *exchange.TaskRequest {
	return &exchange.TaskRequest{Type: exchange.Execute}
}

func (s *ShellModule) Execute(connection connection.Connection, request *exchange.TaskRequest, vars map[string]interface{}) *exchange.TaskResponse {
	response := exchange.NewTaskResponse()
	cmd_parsed, err := templar.GetParsedTemplate(s.Cmd, vars)
	if err != nil {
		response.Error = err
		return response
	}
	response = connection.Execute(cmd_parsed)
	return response
}
