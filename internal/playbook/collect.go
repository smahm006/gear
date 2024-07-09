package playbook

import (
	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/state"
)

func collectHosts(state *state.RunState, play *Play) (map[string]*inventory.Host, error) {
	var err error
	collective_hosts := make(map[string]*inventory.Host)
	limited := len(state.ParsedFlags.Limit) != 0
	if !limited {
		for _, group_name := range play.Groups {
			hosts := state.Inventory.GroupHostsMembership.GroupToHosts[group_name]
			for _, host_name := range hosts {
				host, _ := state.Inventory.GetHost(host_name)
				collective_hosts[host_name] = host
			}
		}
		return collective_hosts, nil
	} else {
		collective_hosts, err = getHostsGivenLimit(state.ParsedFlags.Limit, state, play)
		if err != nil {
			return nil, err
		}
	}
	return collective_hosts, nil
}
