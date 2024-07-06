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

func validateGroups(inventory *inventory.Inventory, play *Play) error {
	for _, group_name := range play.Groups {
		if group_name == "all" {
			var groups []string
			for k := range inventory.Groups {
				groups = append(groups, k)
			}
			play.Groups = groups
			break
		} else {
			_, ok := inventory.GetGroup(group_name)
			if !ok {
				return &PlayValidationError{Name: play.Name, Err: fmt.Errorf("at least one referenced group (%s) is not found in inventory", group_name)}
			}
		}
	}

	return nil
}

func validateHosts(inventory *inventory.Inventory, play *Play) error {
	return nil
}
