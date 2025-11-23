package summarizer

import "github.com/eikendev/cheergo/internal/diff"

// Config captures the data summarizers need to build notifications.
type Config struct {
	GitHubUser string
	LLMApiKey  string
	LLMBaseURL string
	LLMModel   string
}

// Summarizer generates human-readable notification messages from repository diffs.
type Summarizer interface {
	GenerateNotificationMessage(jar *diff.Jar, cfg Config) (string, error)
}
