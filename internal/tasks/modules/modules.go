package modules

import (
	"github.com/smahm006/gear/internal/tasks/requonse"
)

type Module interface {
	Query() requonse.TaskRequest
	Run()
}
