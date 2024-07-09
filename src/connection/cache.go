package connection

import "github.com/smahm006/gear/src/inventory"

type ConnectionCache struct {
	connections map[string]Connection
}

func NewConnectionCache() *ConnectionCache {
	return &ConnectionCache{
		connections: make(map[string]Connection),
	}
}

func (cache *ConnectionCache) GetConnection(host *inventory.Host) (Connection, error) {
	var conn Connection
	conn, ok := cache.connections[host.Name]
	if !ok {
		if host.IsLocal() {
			conn = NewLocalConnection(host)
			cache.connections[host.Name] = conn

		} else {
			conn = NewSshConnection(host)
			cache.connections[host.Name] = conn
		}
	}
	return conn, nil
}
