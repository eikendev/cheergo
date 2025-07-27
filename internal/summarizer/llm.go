// Package summarizer provides interfaces and implementations for generating human-readable summaries of repository changes.
package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/options"
)

type llmNotification struct {
	Message string `json:"message"`
}

// LLMSummarizer implements the Summarizer interface using a large language model.
type LLMSummarizer struct{}

// NewLLMSummarizer returns a new Summarizer that uses a large language model.
func NewLLMSummarizer() Summarizer {
	return &LLMSummarizer{}
}

const (
	systemMessage      = "You are an assistant that generates short, motivating notification messages for developers based on recent GitHub repository activity."
	userPromptTemplate = `You will receive a list of repository changes. Each entry may contain changes in the number of stars, watchers, or subscribers. Your task is to write a single concise and naturally worded plain-text message that highlights the most interesting changes across all repositories.

Requirements:
- The message must be plain text (no Markdown or formatting)
- Use appropriate emojis (e.g., ⭐️, 👀, 🔔, 🚀, 🎉) to enhance tone
- Refer to repository names exactly as given (do not alter them)
- Mention **every repository** included in the input, even if it had only small changes
- Repositories with bigger or more interesting changes should be emphasized
- Repositories with smaller changes can be mentioned more briefly or in passing
- You may reference multiple repositories in the same sentence
- You can be creative, but the message must be relevant to the changes
- The tone should be light, friendly, and motivating
- The message should sound like it was written by a human, not generated
- Avoid rigid, templated language
- Keep the message concise

Here is the user's name: %s

Here is the current time (ISO 8601, UTC): %s

Here is the list of repository changes:
%s`
)

// GenerateNotificationMessage generates a notification message using the LLM.
func (s *LLMSummarizer) GenerateNotificationMessage(jar *diff.Jar, opts *options.Options) (string, error) {
	diffsJSON, err := json.MarshalIndent(jar.Diffs, "", "  ")
	if err != nil {
		slog.Warn("Failed to marshal diffs to JSON", "error", err)
		return "", err
	}

	prompt := fmt.Sprintf(
		userPromptTemplate,
		opts.GitHubUser,
		time.Now().UTC().Format(time.RFC3339),
		string(diffsJSON),
	)

	schema, err := jsonschema.GenerateSchemaForType(llmNotification{})
	if err != nil {
		slog.Warn("Failed to generate JSON schema", "error", err)
		return "", err
	}

	config := openai.DefaultConfig(opts.LLMApiKey)
	config.BaseURL = opts.LLMBaseURL
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: opts.LLMModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "notification",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		slog.Warn("OpenAI API call failed", "error", err)
		return "", err
	}

	var out llmNotification
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &out)
	if err != nil {
		slog.Warn("Failed to unmarshal LLM response", "error", err)
		return "", err
	}

	return out.Message, nil
}
