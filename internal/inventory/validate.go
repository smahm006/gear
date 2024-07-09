package inventory

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/smahm006/gear/internal/utils"
)

type InventoryValidationError struct {
	Path string
	Err  error
}

func (iv *InventoryValidationError) Error() string {
	return fmt.Sprintf("validating %q:\n%v", iv.Path, iv.Err)
}

func validateInventoryPath(path string) ([]byte, error) {
	var yaml []byte
	var err error
	inventory_error := &InventoryValidationError{Path: path}
	if len(path) == 0 {
		yaml, err = utils.ReadFile("inventory.yml")
		if err != nil {
			inventory_error.Err = err
			return nil, inventory_error
		}
		return yaml, nil
	} else {
		yaml, err = utils.ReadFile(path)
		if err != nil {
			inventory_error.Err = err
			return nil, inventory_error
		}
		return yaml, nil
	}
}

func validateInventoryData(path string, i *Inventory) error {
	return_err := &InventoryValidationError{Path: path}
	if len(i.Groups) == 0 {
		return_err.Err = errors.New("no groups found")
		return return_err
	}
	if len(i.Hosts) == 0 {
		return_err.Err = errors.New("no hosts found")
		return return_err
	}

	return nil
}

func validateInventoryValueType(path string, key string, value interface{}, expected_type reflect.Type) error {
	if reflect.TypeOf(value) == expected_type {
		return nil
	}
	return &InventoryValidationError{
		Path: path,
		Err:  fmt.Errorf("invalid yaml format for key %q", key),
	}
}
