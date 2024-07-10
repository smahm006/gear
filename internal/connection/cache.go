package connection

import (
	"github.com/smahm006/gear/internal/inventory"
)

type ConnectionCache struct {
	Connections map[string]Connection
}

func NewConnectionCache() *ConnectionCache {
	return &ConnectionCache{
		Connections: make(map[string]Connection),
	}
}

func (cache *ConnectionCache) GetConnection(host *inventory.Host) (Connection, error) {
	var conn Connection
	conn, ok := cache.Connections[host.Name]
	if !ok {
		if host.IsLocal() {
			conn = NewLocalConnection(host)
			cache.Connections[host.Name] = conn

		} else {
			conn = NewSshConnection(host)
			cache.Connections[host.Name] = conn
		}
	}
	return conn, nil
}
