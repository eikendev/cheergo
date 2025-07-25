// Package main provides the main function as a starting point of this tool.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/containrrr/shoutrrr"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/github"
	"github.com/eikendev/cheergo/internal/options"
	"github.com/eikendev/cheergo/internal/storage"
)

var opts = options.Options{}

func main() {
	kong.Parse(
		&opts,
		kong.Description("Monitors your GitHub repositories for new stars and followers, sending notifications when changes are detected."),
	)

	sender, err := shoutrrr.CreateSender(opts.ShoutrrrURL)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		os.Exit(1)
	}

	data, err := storage.Read(opts.Storage)
	if err != nil {
		slog.Error("Failed to read storage", "error", err)
		os.Exit(1)
	}
	if data == nil {
		slog.Error("Storage data is nil")
		return
	}

	newRepos, err := github.GetRepositories(opts.GitHubUser)
	if err != nil {
		slog.Error("Failed to fetch repositories", "error", err)
		os.Exit(1)
	}

	slog.Info("Fetched repositories",
		"user", opts.GitHubUser,
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

	err = storage.Write(opts.Storage, data)
	if err != nil {
		slog.Error("Failed to write storage", "error", err)
		os.Exit(1)
	}
}
