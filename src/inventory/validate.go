package inventory

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/smahm006/gear/lib/io"
)

type InventoryValidationError struct {
	Path string
	Err  error
}

func (iv *InventoryValidationError) Error() string {
	return fmt.Sprintf("validating %q: %v", iv.Path, iv.Err)
}

func validateInventoryPath(path string) (string, error) {
	if len(path) == 0 {
		path1 := "inventory.yaml"
		path2 := "inventory.yml"
		_, err1 := io.OpenFile("inventory.yaml")
		_, err2 := io.OpenFile("inventory.yml")
		if errors.Is(err1, os.ErrNotExist) || errors.Is(err2, os.ErrNotExist) {
			return "", fmt.Errorf("no inventory file provided or found")
		} else if !errors.Is(err1, os.ErrNotExist) {
			return path1, nil
		} else if !errors.Is(err1, os.ErrNotExist) {
			return path2, nil
		}
	}
	return path, nil
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
