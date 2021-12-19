package github

import (
	"context"

	"github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
)

func GetRepositories(user string) ([]*github.Repository, error) {
	all_repos := []*github.Repository{}
	client := github.NewClient(nil)
	ctx := context.Background()
	n := 100
	page := 1

	for {
		opt := &github.RepositoryListOptions{
			Type:        "owner",
			ListOptions: github.ListOptions{PerPage: n, Page: page},
		}

		new_repos, _, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			return nil, err
		}

		if len(new_repos) > n {
			log.WithFields(log.Fields{
				"page":           page,
				"expected (max)": n,
				"have":           len(new_repos),
			}).Warn("API returned wrong number of repositories for page")
		}

		all_repos = append(all_repos, new_repos...)

		if len(new_repos) != n {
			break
		}

		page++
	}

	return all_repos, nil
}
