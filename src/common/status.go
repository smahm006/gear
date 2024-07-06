package common

import (
	"github.com/smahm006/gear/src/connection"
	"github.com/smahm006/gear/src/inventory"
)

// Status - The situation at a particular time during a process. ‘an update on the status of the bill’
type RunStatus struct {
	ConnectionCache *connection.ConnectionCache
	Hosts           map[string]*inventory.Host
	FailedHosts     map[string]*inventory.Host
	Variables       map[string]interface{}
}

func NewRunStatus(hosts map[string]*inventory.Host, vars map[string]interface{}) *RunStatus {
	return &RunStatus{
		ConnectionCache: connection.NewConnectionCache(),
	}
}
