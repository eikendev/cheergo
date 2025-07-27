// Package diff provides functionality to analyze updates.
package diff

import (
	"log/slog"

	g "github.com/google/go-github/v74/github"

	"github.com/eikendev/cheergo/internal/repository"
)

// Diff holds the difference between the latest update and the previous one.
type Diff struct {
	Stargazers  int
	Subscribers int
	Forks       int
}

// Jar holds a collection of Diff objects and information to notify the user.
type Jar struct {
	Diffs map[string]Diff
}

// NewJar creates a new Jar given information to notify the user.
func NewJar() *Jar {
	return &Jar{
		Diffs: make(map[string]Diff),
	}
}

// Add adds a new Diff into the Jar if a difference in the latest update was detected.
func (d *Jar) Add(name string, is *g.Repository, was *g.Repository) {
	starDiff := is.GetStargazersCount() - was.GetStargazersCount()
	subscribersDiff := is.GetSubscribersCount() - was.GetSubscribersCount()
	forksDiff := is.GetForksCount() - was.GetForksCount()

	slog.Debug("Comparing repository status",
		"starDiff", starDiff,
		"subscribersDiff", subscribersDiff,
		"forksDiff", forksDiff,
		"repository", name,
	)

	if starDiff > 0 || subscribersDiff > 0 || forksDiff > 0 {
		d.Diffs[name] = Diff{
			Stargazers:  starDiff,
			Subscribers: subscribersDiff,
			Forks:       forksDiff,
		}
	}
}

// ComputeDiffs compares newRepos with prevRepos and populates Diffs.
func (d *Jar) ComputeDiffs(newRepos []*g.Repository, prevRepos map[string]g.Repository) {
	for _, is := range newRepos {
		name, err := repository.RepoFullName(is)
		if err != nil {
			slog.Warn("Skipping repository due to missing name", "err", err)
			continue
		}

		was, ok := prevRepos[name]
		if ok {
			d.Add(name, is, &was)
		} else {
			slog.Info("Repository not found in previous data", "repository", name)
		}
	}
}
