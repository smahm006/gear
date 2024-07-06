package playbook

import (
	"fmt"

	"github.com/smahm006/gear/src/inventory"
)

type PlayValidationError struct {
	Name string
	Err  error
}

func (pv *PlayValidationError) Error() string {
	return fmt.Sprintf("validating play %q:\n%v", pv.Name, pv.Err)
}

func validateGroups(inventory *inventory.Inventory, play Play) error {
	for _, group_name := range play.Groups {
		_, ok := inventory.GetGroup(group_name)
		if !ok {
			return &PlayValidationError{Name: play.Name, Err: fmt.Errorf("at least one referenced group (%s) is not found in inventory", group_name)}
		}
	}

	return nil
}

func validateHosts(inventory *inventory.Inventory, play Play) error {
	switch groups := play.Groups.(type) {
	case string:
		group, _ := inventory.GetGroup(groups)
		if len(group.Hosts) == 0 {
			return &PlayValidationError{Name: play.Name, Err: fmt.Errorf("at least one referenced group (%s) is not found in inventory", groups)}
		}
	case []string:
		for _, group := range groups {
			_, ok := inventory.GetGroup(group)
			if !ok {
				return &PlayValidationError{Name: play.Name, Err: fmt.Errorf("at least one referenced group (%s) is not found in inventory", group)}
			}
		}

	}
	return nil
}
