package state

import (
	"github.com/smahm006/gear/internal/connection"
	"github.com/smahm006/gear/internal/inventory"
)

// Status - The situation at a particular time during a process. ‘an update on the status of the bill’
type RunStatus struct {
	ConnectionCache *connection.ConnectionCache
	Hosts           map[string]*inventory.Host
	FailedHosts     map[string]*inventory.Host
}

func NewRunStatus(hosts map[string]*inventory.Host) *RunStatus {
	return &RunStatus{
		ConnectionCache: connection.NewConnectionCache(),
		Hosts:           hosts,
	}
}