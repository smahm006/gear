package connection

import (
	"github.com/smahm006/gear/internal/tasks/requonse"
)

const default_exit_code = 1

type Connection interface {
	Connect() error
	Close() error
	WhoAmI() (string, error)
	Execute(string) *requonse.TaskResponse
	CopyFile(string, string) error
	WriteData(string, string) error
}
