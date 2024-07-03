package connection

import (
	"fmt"

	"github.com/smahm006/gear/src/inventory"
)

type Connection interface {
	Connect() error
	Close() error
	WhoAmI() (string, error)
	Execute() error
	CopyFile() error
	WriteData() error
}

func GetConnection(host *inventory.Host) (Connection, error) {
	var conn Connection
	var env []string
	for k, v := range host.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	if host.IsLocal() {
		conn = NewLocalConnection(env)
	} else {
		conn = NewSshConnection(host, env)
	}
	return conn, nil
}
