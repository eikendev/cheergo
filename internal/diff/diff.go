package diff

import (
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	g "github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
)

type Diff struct {
	Stargazers int
	Watchers   int
	Forks      int
}

type DiffJar struct {
	Diffs  map[string]Diff
	Sender *router.ServiceRouter
}

func NewDiffJar(sender *router.ServiceRouter) *DiffJar {
	return &DiffJar{
		Diffs:  make(map[string]Diff),
		Sender: sender,
	}
}

func (d *DiffJar) Add(name string, is *g.Repository, was *g.Repository) {
	starDiff := is.GetStargazersCount() - was.GetStargazersCount()
	watchDiff := is.GetWatchersCount() - was.GetWatchersCount()
	forksDiff := is.GetForksCount() - was.GetForksCount()

	log.WithFields(log.Fields{
		"starDiff":  starDiff,
		"watchDiff": watchDiff,
		"forksDiff": forksDiff,
		"repo":      name,
	}).Debug("Comparing repository status")

	if starDiff > 0 || watchDiff > 0 || forksDiff > 0 {
		d.Diffs[name] = Diff{
			Stargazers: starDiff,
			Watchers:   watchDiff,
			Forks:      forksDiff,
		}
	}
}

func (d *DiffJar) Send() error {
	if len(d.Diffs) == 0 {
		log.Info("No updates to send")
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

	log.WithFields(log.Fields{
		"msg": msg,
	}).Info("Sending notification")

	err := d.Sender.Send(msg, &types.Params{})
	if len(err) > 0 && err[0] != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}

	return nil
}
