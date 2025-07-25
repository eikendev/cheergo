// Package diff provides functionality to analyze updates and send notifications.
package diff

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	g "github.com/google/go-github/v74/github"
)

// Diff holds the difference between the latest update and the previous one.
type Diff struct {
	Stargazers int
	Watchers   int
	Forks      int
}

// Jar holds a collection of Diff objects and information to notify the user.
type Jar struct {
	Diffs  map[string]Diff
	Sender *router.ServiceRouter
}

// NewJar creates a new Jar given information to notify the user.
func NewJar(sender *router.ServiceRouter) *Jar {
	return &Jar{
		Diffs:  make(map[string]Diff),
		Sender: sender,
	}
}

// Add adds a new Diff into the Jar if a difference in the latest update was detected.
func (d *Jar) Add(name string, is *g.Repository, was *g.Repository) {
	starDiff := is.GetStargazersCount() - was.GetStargazersCount()
	watchDiff := is.GetWatchersCount() - was.GetWatchersCount()
	forksDiff := is.GetForksCount() - was.GetForksCount()

	slog.Debug("Comparing repository status",
		"starDiff", starDiff,
		"watchDiff", watchDiff,
		"forksDiff", forksDiff,
		"repo", name,
	)

	if starDiff > 0 || watchDiff > 0 || forksDiff > 0 {
		d.Diffs[name] = Diff{
			Stargazers: starDiff,
			Watchers:   watchDiff,
			Forks:      forksDiff,
		}
	}
}

// Send sends the update as a notifiction to the user.
func (d *Jar) Send() error {
	if len(d.Diffs) == 0 {
		slog.Info("No updates to send")
		return nil
	}

	msg := ""

	for name, diff := range d.Diffs {
		var updates []string

		if diff.Stargazers > 0 {
			updates = append(updates, fmt.Sprintf("%d new stargazers", diff.Stargazers))
		}
		if diff.Watchers > 0 {
			updates = append(updates, fmt.Sprintf("%d new watchers", diff.Watchers))
		}
		if diff.Forks > 0 {
			updates = append(updates, fmt.Sprintf("%d new forks", diff.Forks))
		}

		msg += fmt.Sprintf("%s has %s!\n", name, strings.Join(updates, ", "))
	}

	slog.Info("Sending notification", "msg", msg)

	err := d.Sender.Send(msg, &types.Params{})
	if len(err) > 0 && err[0] != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}

	return nil
}
