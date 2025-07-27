// Package repository provides functionality related to a GitHub repository.
package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/go-github/v74/github"
)

// GetRepositories returns a list of repositories from a given user.
func GetRepositories(user string) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}
	client := github.NewClient(nil)
	ctx := context.Background()
	n := 100
	page := 1

	for {
		opt := &github.RepositoryListByUserOptions{
			Type:        "owner",
			ListOptions: github.ListOptions{PerPage: n, Page: page},
		}

		newRepos, _, err := client.Repositories.ListByUser(ctx, user, opt)
		if err != nil {
			return nil, err
		}

		if len(newRepos) > n {
			slog.Warn("API returned wrong number of repositories for page",
				"page", page,
				"expected_max", n,
				"have", len(newRepos),
			)
		}

		allRepos = append(allRepos, newRepos...)

		if len(newRepos) != n {
			break
		}

		page++
	}

	return allRepos, nil
}

// RepoFullName returns the "owner/name" string for a given repository.
func RepoFullName(repo *github.Repository) (string, error) {
	if repo == nil {
		return "", errors.New("repository is nil")
	}
	if repo.Owner == nil || repo.Owner.Login == nil || *repo.Owner.Login == "" {
		return "", errors.New("repository owner is missing")
	}
	if repo.Name == nil || *repo.Name == "" {
		return "", errors.New("repository name is missing")
	}
	return *repo.Owner.Login + "/" + *repo.Name, nil
}
