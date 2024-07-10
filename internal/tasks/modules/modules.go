package modules

import (
	"github.com/smahm006/gear/internal/connection"
	"github.com/smahm006/gear/internal/tasks/exchange"
)

type Module interface {
	Query() *exchange.TaskRequest
	Run(request *exchange.TaskRequest, connection connection.Connection) *exchange.TaskResponse
}
