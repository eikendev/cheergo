package summarizer

import (
	"fmt"
	"strings"

	"github.com/eikendev/cheergo/internal/diff"
)

// StaticSummarizer implements Summarizer with a static message generator.
type StaticSummarizer struct{}

// NewStaticSummarizer returns a new Summarizer that generates static messages.
func NewStaticSummarizer() Summarizer {
	return StaticSummarizer{}
}

// GenerateNotificationMessage generates a static notification message from repository diffs.
func (s StaticSummarizer) GenerateNotificationMessage(jar *diff.Jar, _ Config) (string, error) {
	if len(jar.Diffs) == 0 {
		return "", nil
	}

	var sb strings.Builder

	for name, diff := range jar.Diffs {
		var updates []string
		if diff.Stargazers.Diff > 0 {
			updates = append(updates, fmt.Sprintf("%d new stargazers", diff.Stargazers.Diff))
		}
		if diff.Subscribers.Diff > 0 {
			updates = append(updates, fmt.Sprintf("%d new subscribers", diff.Subscribers.Diff))
		}
		if diff.Forks.Diff > 0 {
			updates = append(updates, fmt.Sprintf("%d new forks", diff.Forks.Diff))
		}
		sb.WriteString(fmt.Sprintf("%s has %s!\n", name, strings.Join(updates, ", ")))
	}

	return sb.String(), nil
}
