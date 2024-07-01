package io

import (
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	contents, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
