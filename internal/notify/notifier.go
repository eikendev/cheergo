// Package notify provides notification delivery mechanisms.
package notify

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"

	"github.com/eikendev/cheergo/internal/diff"
)

// Notifier defines the interface for sending notifications about diffs.
type Notifier interface {
	Notify(diffs map[string]diff.Diff) error
}

// ShoutrrrNotifier implements Notifier using a shoutrrr ServiceRouter.
type ShoutrrrNotifier struct {
	Sender *router.ServiceRouter
}

// NewShoutrrrNotifier creates a new ShoutrrrNotifier.
func NewShoutrrrNotifier(sender *router.ServiceRouter) *ShoutrrrNotifier {
	return &ShoutrrrNotifier{Sender: sender}
}

// Notify sends a notification for the given diffs.
func (n *ShoutrrrNotifier) Notify(diffs map[string]diff.Diff) error {
	if len(diffs) == 0 {
		slog.Info("No updates to send")
		return nil
	}

	var msg strings.Builder

	for name, diff := range diffs {
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
		msg.WriteString(fmt.Sprintf("%s has %s!\n", name, strings.Join(updates, ", ")))
	}

	slog.Info("Sending notification", "msg", msg.String())
	errs := n.Sender.Send(msg.String(), &types.Params{})
	if len(errs) > 0 && errs[0] != nil {
		return fmt.Errorf("unable to send notification: %v", errs)
	}
	return nil
}
