package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/containrrr/shoutrrr"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/github"
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

	jar := diff.NewJar(sender)

	for _, is := range newRepos {
		name := fmt.Sprintf("%s/%s", *is.Owner.Login, *is.Name)
		was, ok := data.Repositories[name]
		if ok {
			jar.Add(name, is, &was)
		}
		data.Repositories[name] = *is
	}

	err = jar.Send()
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
