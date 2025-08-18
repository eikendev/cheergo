package commands

import (
	"log/slog"
	"os"

	"github.com/containrrr/shoutrrr"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/logging"
	"github.com/eikendev/cheergo/internal/notify"
	"github.com/eikendev/cheergo/internal/options"
	"github.com/eikendev/cheergo/internal/repository"
	"github.com/eikendev/cheergo/internal/storage"
	"github.com/eikendev/cheergo/internal/summarizer"
)

// RunCommand represents the run subcommand.
type RunCommand struct {
	*options.Options `embed:""`
}

// Run executes the main logic of the application.
func (cmd *RunCommand) Run(_ *options.Options) error {
	logging.Setup(cmd.Verbose)

	sender, err := shoutrrr.CreateSender(cmd.ShoutrrrURL)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		os.Exit(1)
	}

	var summarizerImpl summarizer.Summarizer
	if cmd.LLMApiKey != "" {
		summarizerImpl = summarizer.NewLLMSummarizer()
		slog.Info("Using LLM summarizer")
	} else {
		summarizerImpl = summarizer.NewStaticSummarizer()
		slog.Info("Using static summarizer")
	}

	data, err := storage.Read(cmd.Storage)
	if err != nil {
		slog.Error("Failed to read storage", "error", err)
		os.Exit(1)
	}
	if data == nil {
		slog.Error("Storage data is nil")
		return nil
	}

	newRepos, err := repository.GetRepositories(cmd.GitHubUser)
	if err != nil {
		slog.Error("Failed to fetch repositories", "error", err)
		os.Exit(1)
	}

	slog.Info("Fetched repositories",
		"user", cmd.GitHubUser,
		"count", len(newRepos),
	)

	jar := diff.NewJar()
	jar.ComputeDiffs(newRepos, data.Repositories)

	// Helper to persist latest repository snapshot.
	persist := func() {
		data.UpdateRepositoriesFromSlice(newRepos)

		if err := storage.Write(cmd.Storage, data); err != nil {
			slog.Error("Failed to write storage", "error", err)
			os.Exit(1)
		}
	}

	if len(jar.Diffs) == 0 {
		// Persist even if there are no diffs in case there is no storage file yet.
		persist()

		slog.Info("No repository changes detected", "diff_count", len(jar.Diffs))
		return nil
	}

	messageText, err := summarizerImpl.GenerateNotificationMessage(jar, cmd.Options)
	if err != nil {
		slog.Error("Failed to generate notification message", "error", err)
		os.Exit(1)
	}

	notifier := notify.NewShoutrrrNotifier(sender)
	if err := notifier.NotifyMessage(messageText); err != nil {
		slog.Error("Failed to send notification", "error", err)
		os.Exit(1)
	}

	persist()

	return nil
}
