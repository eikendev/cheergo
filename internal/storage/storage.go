package storage

import (
	"os"

	"github.com/google/go-github/v37/github"
	"gopkg.in/yaml.v3"
)

type Store struct {
	Repositories map[string]github.Repository `yaml:"repositories"`
}

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
