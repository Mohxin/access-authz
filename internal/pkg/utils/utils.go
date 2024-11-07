package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"syscall"

	"sigs.k8s.io/yaml"
)

var ErrNotFound = errors.New("file/directory not found")

func YAMLToJSON(filepath string) ([]byte, error) {
	f, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return yaml.YAMLToJSON(f)
}

func YAMLUnmarshal[T any](filepath string) (T, error) {
	var result T
	data, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) {
			return result, ErrNotFound
		}
		return result, err
	}

	err = yaml.Unmarshal(data, &result)
	return result, err
}

func ReadFileNames(dirPath string) ([]string, error) {
	return readNames(dirPath, func(entry fs.DirEntry) bool {
		return !entry.IsDir()
	})
}

func ReadDirNames(dirPath string) ([]string, error) {
	return readNames(dirPath, func(entry fs.DirEntry) bool {
		return entry.IsDir()
	})
}

func readNames(dirPath string, filter func(entry fs.DirEntry) bool) ([]string, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to read directory [%s]: %w", dirPath, err)
	}
	defer f.Close()

	list, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range list {
		if filter(entry) {
			names = append(names, path.Join(f.Name(), entry.Name()))
		}
	}

	return names, nil
}
