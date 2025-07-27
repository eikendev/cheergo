package summarizer

import (
	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/options"
)

// Summarizer generates human-readable notification messages from repository diffs.
type Summarizer interface {
	GenerateNotificationMessage(jar *diff.Jar, opts *options.Options) (string, error)
}
