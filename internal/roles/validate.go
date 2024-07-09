package roles

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/smahm006/gear/internal/utils"
)

type RoleValidationError struct {
	Name string
	Err  error
}

func (rv *RoleValidationError) Error() string {
	name := strings.Replace(rv.Name, ".yml", "", -1)
	return fmt.Sprintf("validating role %q:\n%v", name, rv.Err)
}

func validateRole(r *Role) ([]byte, error) {
	role_err := &RoleValidationError{Name: r.Name}
	wd, err := os.Getwd()
	if err != nil {
		role_err.Err = err
		return nil, role_err
	}
	// Every role must have a role.yml
	role_dir := fmt.Sprintf("%s/examples/roles/%s", wd, r.Name)
	role_path := fmt.Sprintf("%s/role.yml", role_dir)
	yaml, err := utils.ReadFile(role_path)
	if err != nil {
		role_err.Err = err
		return nil, role_err
	}
	r.Path = role_dir
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
