package modules

import (
	"github.com/smahm006/gear/internal/tasks/exchange"
)

type Module interface {
	Query() exchange.TaskRequest
	Run()
}
