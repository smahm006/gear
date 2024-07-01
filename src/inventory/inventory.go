package inventory

import (
	"github.com/smahm006/gear/lib/io"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Groups map[string]*Group
}

func NewInventory() *Inventory {
	return &Inventory{
		Groups: make(map[string]*Group),
	}
}

func (i *Inventory) LoadInventory(path string) error {
	var processGroups func(gname string, gdata interface{}) (*Group, error)
	processGroups = func(gname string, gdata interface{}) (*Group, error) {
		group := NewGroup(gname)
		for gkey, gvalue := range gdata.(map[string]interface{}) {
			switch gkey {
			case "vars":
				group.Variables = gvalue.((map[string]interface{}))
			case "env":
				group.Environment = gvalue.((map[string]interface{}))
			case "hosts":
				if group.Hosts == nil {
					group.Hosts = make(map[string]*Host)
				}
				switch v := gvalue.(type) {
				case string:
					host := NewHost(gvalue.(string))
					group.Hosts[host.Name] = host
				case map[string]interface{}:
					for hame, hdata := range v {
						host := NewHost(hame)
						for hkey, hvalue := range hdata.(map[string]interface{}) {
							switch hkey {
							case "vars":
								host.Variables = hvalue.((map[string]interface{}))
							case "env":
								host.Environment = hvalue.((map[string]interface{}))
							}
						}
						group.Hosts[host.Name] = host
					}
				case []interface{}:
					for _, hvalue := range v {
						switch v := hvalue.(type) {
						case string:
							host := NewHost(hvalue.(string))
							group.Hosts[host.Name] = host
						case map[string]interface{}:
							for hame, hdata := range v {
								host := NewHost(hame)
								for hhkey, hhvalue := range hdata.(map[string]interface{}) {
									switch hhkey {
									case "vars":
										host.Variables = hhvalue.((map[string]interface{}))
									case "env":
										host.Environment = hhvalue.((map[string]interface{}))
									}
								}
								group.Hosts[host.Name] = host
							}
						}
					}
				}
			default:
				subgroup, err := processGroups(gkey, gvalue)
				if err != nil {
					return nil, err
				}
				if group.SubGroup == nil {
					group.SubGroup = make(map[string]*Group)
				}
				subgroup.ParentGroup = group.Name
				group.SubGroup[gkey] = subgroup
			}
		}
		return group, nil
	}
	yaml_data, err := io.ReadFile(path)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	err = yaml.Unmarshal(yaml_data, &m)
	if err != nil {
		return err
	}
	for gname, gdata := range m {
		group, err := processGroups(gname, gdata)
		if err != nil {
			return err
		}
		i.Groups[gname] = group
	}
	return nil
}
