package commands

import (
	"log/slog"
	"os"

	"github.com/containrrr/shoutrrr"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/github"
	"github.com/eikendev/cheergo/internal/notify"
	"github.com/eikendev/cheergo/internal/options"
	"github.com/eikendev/cheergo/internal/storage"
)

// RunCommand represents the run subcommand.
type RunCommand struct {
	*options.Options `embed:""`
}

// Run executes the main logic of the application.
func (cmd *RunCommand) Run(_ *options.Options) error {
	sender, err := shoutrrr.CreateSender(cmd.ShoutrrrURL)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		os.Exit(1)
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

	newRepos, err := github.GetRepositories(cmd.GitHubUser)
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

	for _, is := range newRepos {
		if is.Owner == nil || is.Name == nil {
			continue
		}
		name := *is.Owner.Login + "/" + *is.Name
		data.Repositories[name] = *is
	}

	notifier := notify.NewShoutrrrNotifier(sender)
	err = notifier.Notify(jar.Diffs)
	if err != nil {
		slog.Error("Failed to send notifications", "error", err)
		os.Exit(1)
	}

	err = storage.Write(cmd.Storage, data)
	if err != nil {
		slog.Error("Failed to write storage", "error", err)
		os.Exit(1)
	}

	return nil
}
