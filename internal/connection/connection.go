package connection

import "github.com/smahm006/gear/internal/tasks"

type Connection interface {
	Connect() error
	Close() error
	WhoAmI() (string, error)
	Execute(string) (tasks.Response, error)
	CopyFile(string, string) error
	WriteData(string, string) error
}
