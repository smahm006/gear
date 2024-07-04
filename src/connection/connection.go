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
	CopyFile(string, string) error
	WriteData(string, string) error
}

func GetConnection(host *inventory.Host) (Connection, error) {
	var conn Connection
	if host.IsLocal() {
		conn = NewLocalConnection(host)
	} else {
		conn = NewSshConnection(host)
	}
	return conn, nil
}
