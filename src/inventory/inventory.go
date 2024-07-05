package inventory

import (
	"fmt"
	"reflect"

	"github.com/smahm006/gear/src/utils"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Groups map[string]*Group
	Hosts  map[string]*Host
}

func NewInventory() *Inventory {
	return &Inventory{
		Groups: make(map[string]*Group),
		Hosts:  make(map[string]*Host),
	}
}

// LoadInventory attempts to unmarshal the inventory file provided.
// The inventory file should be customizable to the point one can have hosts as
// a single string, a list of strings each with the option of having it's own
// variables and environment.
func (i *Inventory) LoadInventory(path string) error {
	path, err := validateInventoryPath(path)
	if err != nil {
		return err
	}
	var processGroups func(gname string, gdata interface{}, parent *Group) (*Group, error)
	processGroups = func(gname string, gdata interface{}, parent *Group) (*Group, error) {
		group := NewGroup(gname)
		if parent != nil {
			group.ParentGroups[parent.Name] = parent
		}
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
				if group.Hosts == nil {
					group.Hosts = make(map[string]*Host)
				}
				switch v := gvalue.(type) {
				case string:
					host := NewHost(gvalue.(string))
					group.Hosts[host.Name] = host
					i.Hosts[host.Name] = host
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
									group.Environment[henvkey] = fmt.Sprint(henvvar)
								}
							}
						}
						group.Hosts[host.Name] = host
						i.Hosts[host.Name] = host
					}
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
											group.Environment[hhenvkey] = fmt.Sprint(hhenvvar)
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
				for _, host := range subgroup.Hosts {
					group.Hosts[host.Name] = host
				}
				group.SubGroups[gkey] = subgroup
			}
		}
		return group, nil
	}
	yaml_data, err := utils.ReadFile(path)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	err = yaml.Unmarshal(yaml_data, &m)
	if err != nil {
		return err
	}
	for gname, gdata := range m {
		group, err := processGroups(gname, gdata, nil)
		if err != nil {
			return err
		}
		i.Groups[gname] = group
	}
	if err = validateInventoryData(path, i); err != nil {
		return err
	}
	return nil
}

func (i *Inventory) GetGroup(name string) (*Group, bool) {
	group, ok := i.Groups[name]
	return group, ok
}

func (i *Inventory) GetHost(name string) (*Host, bool) {
	host, ok := i.Hosts[name]
	return host, ok
}
