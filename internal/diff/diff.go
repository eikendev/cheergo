// Package diff provides functionality to analyze updates.
package diff

import (
	"log/slog"

	g "github.com/google/go-github/v74/github"

	"github.com/eikendev/cheergo/internal/repository"
)

// StatChange holds the absolute and delta values for a repository statistic (e.g., stargazers, subscribers, forks).
type StatChange struct {
	Before int `json:"before"`
	After  int `json:"after"`
	Diff   int `json:"diff"`
}

// Diff holds the before/after/delta values for repository statistics and static metadata for LLM context.
type Diff struct {
	Stargazers  StatChange `json:"stargazers"`
	Subscribers StatChange `json:"subscribers"`
	Forks       StatChange `json:"forks"`
	Description string     `json:"description"`
	Language    string     `json:"language"`
	License     string     `json:"license"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
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

func makeStatChange(before, after int) StatChange {
	return StatChange{
		Before: before,
		After:  after,
		Diff:   after - before,
	}
}

func extractRepoMetadata(repo *g.Repository) (desc, lang, license, created, updated string) {
	desc = repo.GetDescription()
	lang = repo.GetLanguage()
	if repo.GetLicense() != nil {
		license = repo.GetLicense().GetName()
	}
	created = repo.GetCreatedAt().String()
	updated = repo.GetUpdatedAt().String()
	return
}

// Add adds a new Diff into the Jar if a difference in the latest update was detected.
func (d *Jar) Add(name string, is, was *g.Repository) {
	stargazers := makeStatChange(was.GetStargazersCount(), is.GetStargazersCount())
	subscribers := makeStatChange(was.GetSubscribersCount(), is.GetSubscribersCount())
	forks := makeStatChange(was.GetForksCount(), is.GetForksCount())

	slog.Debug("Comparing repository status",
		"starDiff", stargazers.Diff,
		"subscribersDiff", subscribers.Diff,
		"forksDiff", forks.Diff,
		"repository", name,
	)

	if stargazers.Diff > 0 || subscribers.Diff > 0 || forks.Diff > 0 {
		slog.Info("Repository change detected",
			"repository", name,
			"stargazers", stargazers,
			"subscribers", subscribers,
			"forks", forks,
		)

		desc, lang, license, created, updated := extractRepoMetadata(is)
		d.Diffs[name] = Diff{
			Stargazers:  stargazers,
			Subscribers: subscribers,
			Forks:       forks,
			Description: desc,
			Language:    lang,
			License:     license,
			CreatedAt:   created,
			UpdatedAt:   updated,
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
