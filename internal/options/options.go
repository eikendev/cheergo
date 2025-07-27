// Package options defines the global options of this tool.
package options

// Options represents the global options of this tool.
type Options struct {
	Storage     string `name:"storage" help:"The storage file." type:"file" default:"storage.yml" env:"CHEERGO_STORAGE"`
	ShoutrrrURL string `name:"shoutrrr-url" help:"The URL for Shoutrrr." required:"true" env:"CHEERGO_SHOUTRRR_URL"`
	GitHubUser  string `name:"github-user" help:"The name of the user to monitor." required:"true" env:"CHEERGO_GITHUB_USER"`
	LLMApiKey   string `name:"llm-api-key" help:"API key for LLM (OpenRouter/OpenAI-compatible). If not set, static notifications are used." env:"CHEERGO_LLM_API_KEY"`
	LLMBaseURL  string `name:"llm-base-url" help:"Base URL for LLM API." default:"https://openrouter.ai/api/v1" env:"CHEERGO_LLM_BASE_URL"`
	LLMModel    string `name:"llm-model" help:"LLM model to use." default:"google/gemini-2.5-flash-lite-preview-06-17" env:"CHEERGO_LLM_MODEL"`
	Verbose     bool   `name:"verbose" help:"Enable verbose (debug) logging." env:"CHEERGO_VERBOSE"`
}
