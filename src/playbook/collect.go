package playbook

import (
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/inventory"
)

func collectHosts(state *common.RunState, play Play) map[string]*inventory.Host {
	var collective_hosts map[string]*inventory.Host
	limited := len(state.ParsedFlags.Limit) != 0
	if !limited {
		for _, group_name := range play.Groups {
			group, _ := state.Inventory.GetGroup(group_name)
			for k, v := range group.Hosts {
				collective_hosts[k] = v
			}
		}
	} else {

	}
	return collective_hosts
}

func collectVars(state *common.RunState, play Play) map[string]interface{} {
	var collective_vars map[string]interface{}
	// Start with inventory vars
	return collective_vars
}
