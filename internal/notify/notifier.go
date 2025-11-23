// Package notify provides notification delivery mechanisms.
package notify

import (
	"errors"
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
	var sendErr error
	for _, err := range errs {
		if err != nil {
			sendErr = errors.Join(sendErr, err)
		}
	}
	if sendErr != nil {
		return fmt.Errorf("unable to send notification: %w", sendErr)
	}

	slog.Info("Notification sent", "message", msg)

	return nil
}
