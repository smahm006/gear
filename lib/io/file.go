package io

import (
	"os"
	"path/filepath"
)

func OpenFile(path string) (*os.File, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

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
