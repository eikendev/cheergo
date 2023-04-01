// Package options defines the global options of this tool.
package options

// Options represents the global options of this tool.
type Options struct {
	Storage     string `name:"storage" help:"The storage file." type:"file" default:"storage.yml" env:"CHEERGO_STORAGE"`
	ShoutrrrURL string `name:"shoutrrr-url" help:"The URL for Shoutrrr." required:"true" env:"CHEERGO_SHOUTRRR_URL"`
	GitHubUser  string `name:"github-user" help:"The name of the user to monitor." required:"true" env:"CHEERGO_GITHUB_USER"`
}
