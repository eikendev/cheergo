// Package main provides the main function as a starting point of this tool.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/containrrr/shoutrrr"
	log "github.com/sirupsen/logrus"

	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/github"
	"github.com/eikendev/cheergo/internal/options"
	"github.com/eikendev/cheergo/internal/storage"
)

var opts = options.Options{}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

func main() {
	kong.Parse(
		&opts,
		kong.Description(fmt.Sprintf("%s (%s)", version, date)),
	)

	sender, err := shoutrrr.CreateSender(opts.ShoutrrrURL)
	if err != nil {
		log.Fatal(err)
	}

	data, err := storage.Read(opts.Storage)
	if err != nil {
		log.Fatal(err)
	}
	if data == nil {
		log.Fatal("Storage data is nil")
		return
	}

	newRepos, err := github.GetRepositories(opts.GitHubUser)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"user":  opts.GitHubUser,
		"count": len(newRepos),
	}).Info("Fetched repositories")

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
		log.Fatal(err)
	}

	err = storage.Write(opts.Storage, data)
	if err != nil {
		log.Fatal(err)
	}
}
