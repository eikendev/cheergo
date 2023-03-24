package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/containrrr/shoutrrr"
	"github.com/eikendev/cheergo/internal/diff"
	"github.com/eikendev/cheergo/internal/github"
	"github.com/eikendev/cheergo/internal/options"
	"github.com/eikendev/cheergo/internal/storage"

	log "github.com/sirupsen/logrus"
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

	sender, err := shoutrrr.CreateSender(opts.ShoutrrrUrl)
	if err != nil {
		log.Fatal(err)
	}

	data, err := storage.Read(opts.Storage)
	if err != nil {
		log.Fatal(err)
	}

	new_repos, err := github.GetRepositories(opts.GitHubUser)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"user":  opts.GitHubUser,
		"count": len(new_repos),
	}).Info("Fetched repositories")

	jar := diff.NewDiffJar(sender)

	for _, is := range new_repos {
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
