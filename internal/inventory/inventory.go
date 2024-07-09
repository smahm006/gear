package inventory

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

type GroupHostsMembership struct {
	GroupToHosts map[string][]string
	HostsToGroup map[string][]string
}

type Inventory struct {
	Groups               map[string]*Group
	Hosts                map[string]*Host
	GroupHostsMembership *GroupHostsMembership
}

func NewInventory() *Inventory {
	return &Inventory{
		Groups: make(map[string]*Group),
		Hosts:  make(map[string]*Host),
		GroupHostsMembership: &GroupHostsMembership{
			GroupToHosts: make(map[string][]string),
			HostsToGroup: make(map[string][]string),
		},
	}
}

// LoadInventory attempts to unmarshal the inventory file provided.
// The inventory file should be flexible to the point one can have hosts as
// a single string or a list of strings, each with the option of having it's own
// variables and environment.
func (i *Inventory) LoadInventory(path string) error {
	GroupToHosts := i.GroupHostsMembership.GroupToHosts
	HostsToGrops := i.GroupHostsMembership.HostsToGroup
	var processGroups func(gname string, gdata interface{}, parent *Group) (*Group, error)
	processGroups = func(gname string, gdata interface{}, parent *Group) (*Group, error) {
		group := NewGroup(gname)
		group.ParentGroups[parent.Name] = parent
		if err := validateInventoryValueType(path, gname, gdata, reflect.TypeOf(map[string]interface{}{})); err != nil {
			return group, err
		}
		if gname == "local" {
			localhost := NewHost("127.0.0.1")
			group.Hosts["127.0.0.1"] = localhost
			i.Hosts["127.0.0.1"] = localhost
		}

		for gkey, gvalue := range gdata.(map[string]interface{}) {
			switch gkey {
			case "vars":
				if err := validateInventoryValueType(path, gkey, gvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
					return group, err
				}
				group.Variables = gvalue.((map[string]interface{}))
			case "env":
				if err := validateInventoryValueType(path, gkey, gvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
					return group, err
				}
				for genvkey, genvvar := range gvalue.((map[string]interface{})) {
					group.Environment[genvkey] = fmt.Sprint(genvvar)
				}
			case "hosts":
				switch v := gvalue.(type) {
				// hosts: 10.10.10.1
				case string:
					host := NewHost(gvalue.(string))
					group.Hosts[host.Name] = host
					i.Hosts[host.Name] = host
				// hosts:
				//  10.10.10.1:
				//    vars:
				//      key: value
				case map[string]interface{}:
					for hname, hdata := range v {
						host := NewHost(hname)
						if err := validateInventoryValueType(path, hname, hdata, reflect.TypeOf(map[string]interface{}{})); err != nil {
							return group, err
						}
						for hkey, hvalue := range hdata.(map[string]interface{}) {
							switch hkey {
							case "vars":
								if err := validateInventoryValueType(path, hkey, hvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
									return group, err
								}
								host.Variables = hvalue.((map[string]interface{}))
							case "env":
								if err := validateInventoryValueType(path, hkey, hvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
									return group, err
								}
								for henvkey, henvvar := range hvalue.((map[string]interface{})) {
									host.Environment[henvkey] = fmt.Sprint(henvvar)
								}
							}
						}
						group.Hosts[host.Name] = host
						i.Hosts[host.Name] = host
					}
				// hosts:
				//  - 10.10.10.1
				//  - 10.10.10.2
				//  - 10.10.10.3:
				//      vars:
				//        key: value
				case []interface{}:
					for _, hvalue := range v {
						switch v := hvalue.(type) {
						case string:
							host := NewHost(hvalue.(string))
							group.Hosts[host.Name] = host
							i.Hosts[host.Name] = host
						case map[string]interface{}:
							for hname, hdata := range v {
								host := NewHost(hname)
								if err := validateInventoryValueType(path, hname, hdata, reflect.TypeOf(map[string]interface{}{})); err != nil {
									return group, err
								}
								for hhkey, hhvalue := range hdata.(map[string]interface{}) {
									switch hhkey {
									case "vars":
										if err := validateInventoryValueType(path, hhkey, hhvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
											return group, err
										}
										host.Variables = hhvalue.((map[string]interface{}))
									case "env":
										if err := validateInventoryValueType(path, hhkey, hhvalue, reflect.TypeOf(map[string]interface{}{})); err != nil {
											return group, err
										}
										for hhenvkey, hhenvvar := range hhvalue.((map[string]interface{})) {
											host.Environment[hhenvkey] = fmt.Sprint(hhenvvar)
										}
									}
								}
								group.Hosts[host.Name] = host
								i.Hosts[host.Name] = host
							}
						}
					}
				}
			// if key is not hosts, vars or env we assume it is a subgroup
			default:
				subgroup, err := processGroups(gkey, gvalue, group)
				if err != nil {
					return nil, err
				}
				for host_name, host := range subgroup.Hosts {
					group.Hosts[host_name] = host
					GroupToHosts[subgroup.Name] = append(GroupToHosts[subgroup.Name], host_name)
					HostsToGrops[host_name] = append(HostsToGrops[host_name], subgroup.Name)
				}
				group.SubGroups[gkey] = subgroup
			}
		}
		return group, nil
	}
	yaml_data, err := validateInventoryPath(path)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	err = yaml.Unmarshal(yaml_data, &m)
	if err != nil {
		return err
	}
	for gname, gdata := range m {
		group, err := processGroups(gname, gdata, &Group{})
		if err != nil {
			return err
		}
		// Add top level groups
		i.Groups[gname] = group
		for host_name := range group.Hosts {
			GroupToHosts[group.Name] = append(GroupToHosts[group.Name], host_name)
			HostsToGrops[host_name] = append(HostsToGrops[host_name], group.Name)
		}
	}
	if err = validateInventoryData(path, i); err != nil {
		return err
	}
	// Update variable and environment of groups and hosts to include parent group
	// environment variables. Need to do this after un-marshalling is complete as
	// you cannot predict how go traverses a map
	var updateEnvVars func(m map[string]*Group)
	updateEnvVars = func(m map[string]*Group) {
		for _, group := range m {
			if group.SubGroups != nil {
				updateEnvVars(group.SubGroups)
			}
			for key, value := range group.Environment {
				for _, host := range group.Hosts {
					host.Environment[key] = value
				}
			}
			for key, value := range group.Variables {
				for _, host := range group.Hosts {
					host.Variables[key] = value
				}
			}
			for _, parentgroup := range group.ParentGroups {
				for key, value := range parentgroup.Environment {
					group.Environment[key] = value
					for _, host := range group.Hosts {
						host.Environment[key] = value
					}
				}
				for key, value := range parentgroup.Variables {
					group.Variables[key] = value
					for _, host := range group.Hosts {
						host.Variables[key] = value
					}
				}
			}
		}
	}
	updateEnvVars(i.Groups)
	return nil
}

func getGroupNested(m map[string]*Group, name string) (*Group, bool) {
	if group, exists := m[name]; exists {
		return group, exists
	}
	for _, group := range m {
		if group.ParentGroups != nil {
			if group, exists := getGroupNested(group.ParentGroups, name); exists {
				return group, exists
			}
		}
	}
	return nil, false
}

func (i *Inventory) GetGroup(name string) (*Group, bool) {
	return getGroupNested(i.Groups, name)
}

func (i *Inventory) GetHost(name string) (*Host, bool) {
	host, ok := i.Hosts[name]
	return host, ok
}
