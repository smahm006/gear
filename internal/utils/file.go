package utils

import (
	"io"
	"os"
	"path/filepath"
)

func OpenFile(path string) (*os.File, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o660)
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

func CopyFile(src string, dest string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() {
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}
