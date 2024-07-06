package modules

import "github.com/smahm006/gear/src/common"

type Module interface {
	Query() common.TaskRequest
	Run()
}
