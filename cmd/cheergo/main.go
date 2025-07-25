// Package main is the entry point for the cheergo CLI tool.
package main

import (
	"github.com/alecthomas/kong"

	"github.com/eikendev/cheergo/internal/commands"
	"github.com/eikendev/cheergo/internal/options"
)

type CLI struct {
	Run     commands.RunCommand     `cmd:"" help:"Monitor for new stars and followers on your GitHub repositories."`
	Version commands.VersionCommand `cmd:"" help:"Show version information."`
}

func main() {
	var cli CLI
	var opts options.Options
	kctx := kong.Parse(&cli,
		kong.Description("cheergo is a CLI tool for monitoring your GitHub repositories. It can notify you about new stars and followers."),
		kong.UsageOnError(),
		kong.Bind(&opts),
	)

	err := kctx.Run()
	kctx.FatalIfErrorf(err)
}
