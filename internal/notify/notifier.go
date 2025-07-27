// Package notify provides notification delivery mechanisms.
package notify

import (
	"fmt"
	"log/slog"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Notifier defines the interface for sending notifications about diffs.
type Notifier interface {
	NotifyMessage(msg string) error
}

// ShoutrrrNotifier implements Notifier using a shoutrrr ServiceRouter.
type ShoutrrrNotifier struct {
	Sender *router.ServiceRouter
}

// NewShoutrrrNotifier creates a new ShoutrrrNotifier.
func NewShoutrrrNotifier(sender *router.ServiceRouter) *ShoutrrrNotifier {
	return &ShoutrrrNotifier{Sender: sender}
}

// NotifyMessage sends a raw notification message.
func (n *ShoutrrrNotifier) NotifyMessage(msg string) error {
	if msg == "" {
		slog.Info("No message to send")
		return nil
	}

	errs := n.Sender.Send(msg, &types.Params{})
	if len(errs) > 0 && errs[0] != nil {
		return fmt.Errorf("unable to send notification: %v", errs)
	}

	slog.Info("Notification sent", "message", msg)

	return nil
}
