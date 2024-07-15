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
	Variables       map[string]interface{}
}

func NewRunStatus(hosts map[string]*inventory.Host, vars map[string]interface{}) *RunStatus {
	collected_vars := make(map[string]interface{})
	for _, host := range hosts {
		for key, value := range host.Variables {
			collected_vars[key] = value
		}
	}
	for key, value := range vars {
		collected_vars[key] = value
	}
	return &RunStatus{
		ConnectionCache: connection.NewConnectionCache(),
		Hosts:           hosts,
		Variables:       collected_vars,
	}
}
