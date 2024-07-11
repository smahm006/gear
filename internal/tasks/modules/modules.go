package modules

import (
	"github.com/smahm006/gear/internal/connection"
	"github.com/smahm006/gear/internal/tasks/exchange"
)

type Module interface {
	Query() *exchange.TaskRequest
	Run(connection connection.Connection, request *exchange.TaskRequest, with *ModuleWith, and *ModuleAnd) *exchange.TaskResponse
}
