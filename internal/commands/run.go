package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	gh "github.com/google/go-github/v74/github"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/notify"
	"github.com/eikendev/cheergo/internal/repository"
	"github.com/eikendev/cheergo/internal/storage"
	"github.com/eikendev/cheergo/internal/summarizer"
)

const repoFetchTimeout = 30 * time.Second

// RunCommand represents the run subcommand.
type RunCommand struct {
	Storage     string `name:"storage" help:"The storage file." type:"file" default:"storage.yml" env:"CHEERGO_STORAGE"`
	ShoutrrrURL string `name:"shoutrrr-url" help:"The URL for Shoutrrr." required:"true" env:"CHEERGO_SHOUTRRR_URL"`
	GitHubUser  string `name:"github-user" help:"The name of the user to monitor." required:"true" env:"CHEERGO_GITHUB_USER"`
	LLMApiKey   string `name:"llm-api-key" help:"API key for LLM (OpenRouter/OpenAI-compatible). If not set, static notifications are used." env:"CHEERGO_LLM_API_KEY"`
	LLMBaseURL  string `name:"llm-base-url" help:"Base URL for LLM API." default:"https://openrouter.ai/api/v1" env:"CHEERGO_LLM_BASE_URL"`
	LLMModel    string `name:"llm-model" help:"LLM model to use." default:"google/gemini-2.5-flash" env:"CHEERGO_LLM_MODEL"`
}

func initRun(cmd *RunCommand) (*router.ServiceRouter, summarizer.Summarizer, *storage.Store, error) {
	sender, err := shoutrrr.CreateSender(cmd.ShoutrrrURL)
	if err != nil {
		slog.Error("Failed to create sender", "error", err)
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	return sender, summarizerImpl, data, nil
}

func fetchRepositories(user string) ([]*gh.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repoFetchTimeout)
	defer cancel()

	newRepos, err := repository.GetRepositories(ctx, user)
	if err != nil {
		slog.Error("Failed to fetch repositories", "error", err)
		return nil, err
	}

	return newRepos, nil
}

func persistSnapshot(path string, data *storage.Store, repos []*gh.Repository) error {
	data.UpdateRepositoriesFromSlice(repos)

	if err := storage.Write(path, data); err != nil {
		slog.Error("Failed to write storage", "error", err)
		return err
	}

	return nil
}

func notifyChanges(sender *router.ServiceRouter, summarizerImpl summarizer.Summarizer, jar *diff.Jar, cfg summarizer.Config) error {
	messageText, err := summarizerImpl.GenerateNotificationMessage(jar, cfg)
	if err != nil {
		slog.Error("Failed to generate notification message", "error", err)
		return err
	}

	notifier := notify.NewShoutrrrNotifier(sender)
	if err := notifier.NotifyMessage(messageText); err != nil {
		slog.Error("Failed to send notification", "error", err)
		return err
	}

	return nil
}

// Run executes the main logic of the application.
func (cmd *RunCommand) Run() error {
	sender, summarizerImpl, data, err := initRun(cmd)
	if err != nil {
		return err
	}

	newRepos, err := fetchRepositories(cmd.GitHubUser)
	if err != nil {
		return err
	}

	slog.Info("Fetched repositories",
		"user", cmd.GitHubUser,
		"count", len(newRepos),
	)

	jar := diff.NewJar()
	jar.ComputeDiffs(newRepos, data.Repositories)

	if len(jar.Diffs) == 0 {
		// Persist even if there are no diffs in case there is no storage file yet.
		if err := persistSnapshot(cmd.Storage, data, newRepos); err != nil {
			return err
		}

		slog.Info("No repository changes detected", "diff_count", len(jar.Diffs))
		return nil
	}

	if err := notifyChanges(sender, summarizerImpl, jar, summarizer.Config{
		GitHubUser: cmd.GitHubUser,
		LLMApiKey:  cmd.LLMApiKey,
		LLMBaseURL: cmd.LLMBaseURL,
		LLMModel:   cmd.LLMModel,
	}); err != nil {
		return err
	}

	if err := persistSnapshot(cmd.Storage, data, newRepos); err != nil {
		return err
	}

	return nil
}
