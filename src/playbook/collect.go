package playbook

import (
	"github.com/smahm006/gear/src/common"
	"github.com/smahm006/gear/src/inventory"
)

func collectHosts(state *common.RunState, play *Play) (map[string]*inventory.Host, error) {
	collective_groups := make(map[string]*inventory.Group)
	collective_hosts := make(map[string]*inventory.Host)
	limited := len(state.ParsedFlags.Limit) != 0
	for _, group_name := range play.Groups {
		group, _ := state.Inventory.GetGroup(group_name)
		collective_groups[group_name] = group
		for k, v := range group.Hosts {
			collective_hosts[k] = v
		}
	}
	if !limited {
		return collective_hosts, nil
	} else {
		_, err := getHostsGivenLimit(state.ParsedFlags.Limit, collective_groups)
		if err != nil {
			return nil, err
		}
	}
	return collective_hosts, nil
}

func collectVars(state *common.RunState, play *Play) map[string]interface{} {
	collective_vars := make(map[string]interface{})
	// Start with inventory vars
	return collective_vars
}
