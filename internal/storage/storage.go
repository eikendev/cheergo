// Package storage provides convenience functions related to the file system.
package storage

import (
	"log/slog"
	"os"

	gh "github.com/google/go-github/v74/github"
	"gopkg.in/yaml.v3"

	"github.com/eikendev/cheergo/internal/repository"
)

// Store represents a collection of GitHub repositories.
type Store struct {
	Repositories map[string]gh.Repository `yaml:"repositories"`
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
		store.Repositories = make(map[string]gh.Repository)
	}

	return &store, nil
}

// UpdateRepositoriesFromSlice updates the Store's Repositories map from a slice of *gh.Repository.
func (s *Store) UpdateRepositoriesFromSlice(repos []*gh.Repository) {
	if s.Repositories == nil {
		s.Repositories = make(map[string]gh.Repository)
	}

	for _, repo := range repos {
		name, err := repository.RepoFullName(repo)
		if err == nil {
			s.Repositories[name] = *repo
		}
	}

	slog.Debug("Updated repositories in store", "count", len(s.Repositories))
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

	slog.Info("Wrote storage file", "path", path, "repo_count", len(store.Repositories))

	return nil
}
