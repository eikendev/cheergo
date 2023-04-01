// Package github provides functionality related to a GitHub account.
package github

import (
	"context"

	"github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
)

// GetRepositories returns a list of repositories from a given user.
func GetRepositories(user string) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}
	client := github.NewClient(nil)
	ctx := context.Background()
	n := 100
	page := 1

	for {
		opt := &github.RepositoryListOptions{
			Type:        "owner",
			ListOptions: github.ListOptions{PerPage: n, Page: page},
		}

		newRepos, _, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			return nil, err
		}

		if len(newRepos) > n {
			log.WithFields(log.Fields{
				"page":           page,
				"expected (max)": n,
				"have":           len(newRepos),
			}).Warn("API returned wrong number of repositories for page")
		}

		allRepos = append(allRepos, newRepos...)

		if len(newRepos) != n {
			break
		}

		page++
	}

	return allRepos, nil
}
