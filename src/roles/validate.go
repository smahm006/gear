package roles

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/smahm006/gear/src/utils"
)

type RoleValidationError struct {
	Name string
	Err  error
}

func (iv *RoleValidationError) Error() string {
	name := strings.Replace(iv.Name, ".yml", "", -1)
	return fmt.Sprintf("validating role %q:\n%v", name, iv.Err)
}

func validateRole(r *Role) ([]byte, error) {
	role_err := &RoleValidationError{Name: r.Name}
	wd, err := os.Getwd()
	if err != nil {
		role_err.Err = err
		return nil, role_err
	}
	// Every role must have a role.yml
	path := fmt.Sprintf("%s/examples/roles/%s/role.yml", wd, r.Name)
	yaml, err := utils.ReadFile(path)
	if err != nil {
		role_err.Err = err
		return nil, role_err
	}
	r.Path = path
	return yaml, nil
}

func validateRoleData(i *Role) error {
	return_err := &RoleValidationError{Name: i.Name}
	if len(i.Tasks) == 0 {
		return_err.Err = errors.New("no tasks found")
		return return_err
	}
	return nil
}
