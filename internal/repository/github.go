// Package repository provides functionality related to a GitHub repository.
package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/go-github/v74/github"
)

const perPage = 100

// GetRepositories returns a list of repositories from a given user.
func GetRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	allRepos := []*github.Repository{}
	client := github.NewClient(nil)
	opt := &github.RepositoryListByUserOptions{
		Type:        "owner",
		ListOptions: github.ListOptions{PerPage: perPage},
	}

	for {
		newRepos, resp, err := client.Repositories.ListByUser(ctx, user, opt)
		if err != nil {
			return nil, err
		}

		if len(newRepos) > perPage {
			slog.Warn("API returned unexpected number of repositories for page",
				"page", opt.Page,
				"expected_max", perPage,
				"have", len(newRepos),
			)
		}

		allRepos = append(allRepos, newRepos...)

		if resp == nil || resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
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
