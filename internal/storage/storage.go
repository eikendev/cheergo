// Package storage provides convenience functions related to the file system.
package storage

import (
	"os"

	"github.com/google/go-github/v74/github"
	"gopkg.in/yaml.v3"
)

// Store represents a collection of GitHub repositories.
type Store struct {
	Repositories map[string]github.Repository `yaml:"repositories"`
}

// Read returns a deserialized Store object from a given file.
func Read(path string) (*Store, error) {
	yfile, err := os.ReadFile(path) //#nosec G304
	if os.IsNotExist(err) {
		yfile = []byte{}
	} else if err != nil {
		return nil, err
	}

	var store Store

	err = yaml.Unmarshal(yfile, &store)
	if err != nil {
		return nil, err
	}

	if store.Repositories == nil {
		store.Repositories = make(map[string]github.Repository)
	}

	return &store, nil
}

// Write serializes and writes a Store object to a given file.
func Write(path string, store *Store) error {
	out, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, out, 0o600)
	if err != nil {
		return err
	}

	return nil
}
