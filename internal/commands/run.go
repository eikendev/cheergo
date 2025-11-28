package commands

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/containrrr/shoutrrr"

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

// Run executes the main logic of the application.
func (cmd *RunCommand) Run() error {
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

	ctx, cancel := context.WithTimeout(context.Background(), repoFetchTimeout)
	newRepos, err := repository.GetRepositories(ctx, cmd.GitHubUser)
	cancel()
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

	messageText, err := summarizerImpl.GenerateNotificationMessage(jar, summarizer.Config{
		GitHubUser: cmd.GitHubUser,
		LLMApiKey:  cmd.LLMApiKey,
		LLMBaseURL: cmd.LLMBaseURL,
		LLMModel:   cmd.LLMModel,
	})
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
